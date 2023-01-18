package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/siddontang/go/ioutil2"
	"io"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/request"
	"metalflow/pkg/response"
	"metalflow/pkg/service"
	"metalflow/pkg/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// UploadFileChunkExists determine whether a file block exists.
func UploadFileChunkExists(c *gin.Context) {
	var filePart request.FilePartInfo
	_ = c.ShouldBind(&filePart)
	// verification request.
	err := filePart.ValidateReq()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	filePart.Complete, filePart.Uploaded = getUploadedChunkNumbers(&filePart)
	response.SuccessWithData(filePart)
}

// UploadMerge merge shard files.
// nolint:funlen
func UploadMerge(c *gin.Context) { //nolint:gocyclo
	var fileMerge request.FileMergeInfo
	_ = c.ShouldBind(&fileMerge)

	addressIds := fileMerge.AddressIds
	filePart := fileMerge.FilePartInfo
	rootDir := filePart.GetUploadRootPath()
	mergeFileName := fmt.Sprintf("%s/%s", rootDir, filePart.Filename)
	// create a merge file.
	mergeFile, err := os.OpenFile(mergeFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	defer func(mergeFile *os.File) {
		_ = mergeFile.Close()
	}(mergeFile)

	totalChunk := int(filePart.GetTotalChunk())
	chunkSize := int(filePart.ChunkSize)
	var chunkNumbers []int
	for i := 0; i < totalChunk; i++ {
		chunkNumbers = append(chunkNumbers, i+1)
	}

	// enable go routines to merge files concurrently.
	// If the total number of file blocks is too large, the performance will decrease instead,
	// so you need to configure an appropriate number of coroutines.
	var count = int(global.Conf.Upload.MergeConcurrentCount)
	chunkCount := len(chunkNumbers) / count
	// The last group is considered exactly divisible by default
	lastChunkCount := chunkCount
	if len(chunkNumbers)%count > 0 || count == 1 {
		lastChunkCount = len(chunkNumbers)%count + chunkCount
	}
	// convert to a two-dimensional array, and each set of data is assigned to a routine for use.
	chunks := make([][]int, count)
	for i := 0; i < count; i++ {
		if i < count-1 {
			chunks[i] = chunkNumbers[i*chunkCount : (i+1)*chunkCount]
		} else {
			chunks[i] = chunkNumbers[i*chunkCount : i*chunkCount+lastChunkCount]
		}
	}
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(arr []int) {
			defer wg.Done()
			for _, item := range arr {
				func() {
					currentChunkName := filePart.GetChunkFilename(uint(item))
					exists := ioutil2.FileExists(currentChunkName)
					if exists {
						// read file fragments.
						f, err := os.OpenFile(currentChunkName, os.O_RDONLY, os.ModePerm) //nolint:govet
						if err != nil {
							response.FailWithMsg(err.Error())
							return
						}
						defer func() {
							// close file.
							_ = f.Close()
						}()
						b, err := io.ReadAll(f)
						if err != nil {
							response.FailWithMsg(err.Error())
							return
						}
						// start writing from the specified position.
						_, _ = mergeFile.WriteAt(b, int64((item-1)*chunkSize))
					}
				}()
			}
		}(chunks[i])
	}
	// wait for all coroutines to finish processing.
	wg.Wait()

	ids := utils.Str2UintArr(addressIds)
	m := &MergeInfo{path: mergeFileName}

	fileMetric := grpc.FileMetric{
		FilePath:   mergeFileName,
		RemoteDir:  fileMerge.RemoteDir,
		IsRunnable: fileMerge.Runnable,
		FileGetter: m,
	}
	s := service.New(c)
	err = s.BatchUploadByIds(fileMetric, ids)
	// after the file is transferred to the corresponding machine, delete the path where the fragmented file is located and the original file.
	_ = os.RemoveAll(filePart.GetChunkRootPath())
	_ = os.Remove(mergeFileName)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// write back file information.
	var res response.UploadMergeResponseStruct
	res.Output = ""
	response.SuccessWithData(res)
}

type MergeInfo struct {
	path string
}

// GetFile get file object.
func (m *MergeInfo) GetFile() (io.ReadCloser, error) {
	file, err := os.Open(m.path)
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("fail to open the file：%s", m.path))
		return nil, err
	}
	return file, nil
}

// UploadFile upload files (small files are directly a single file, if a very large file may be a single fragment).
func UploadFile(c *gin.Context) {
	// limit file maximum memory. Binary shift xxxMB.
	err := c.Request.ParseMultipartForm(int64(global.Conf.Upload.SingleMaxSize) << 20) //nolint:gomnd
	if err != nil {
		response.FailWithMsg(fmt.Sprintf("the file size exceeds the maximum %dMB", global.Conf.Upload.SingleMaxSize))
		return
	}
	// read file fragments.
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMsg("unable to read file")
		return
	}

	// read file fragmentation parameters.
	var filePart request.FilePartInfo
	// current size.
	currentSize := uint(header.Size)
	filePart.CurrentSize = &currentSize
	// chunk number.
	filePart.ChunkNumber = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("chunkNumber")))
	// chunk size.
	filePart.ChunkSize = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("chunkSize")))
	// total size.
	filePart.TotalSize = utils.Str2Uint(strings.TrimSpace(c.Request.FormValue("totalSize")))
	// uniquely identifies.
	filePart.Identifier = strings.TrimSpace(c.Request.FormValue("identifier"))
	// file name.
	filePart.Filename = strings.TrimSpace(c.Request.FormValue("filename"))

	// request validate.
	err = filePart.ValidateReq()
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// get chunk filename.
	chunkName := filePart.GetChunkFilename(filePart.ChunkNumber)
	// create a folder that does not exist.
	chunkDir, _ := filepath.Split(chunkName)
	err = os.MkdirAll(chunkDir, os.ModePerm)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// save chunk file.
	out, err := os.Create(chunkName)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	// copy the contents of file to out.
	_, err = io.Copy(out, file)
	if err != nil {
		response.FailWithMsg(err.Error())
		return
	}

	// check file chunk integrity.
	filePart.CurrentCheckChunkNumber = 1
	filePart.Complete = checkChunkComplete(&filePart)
	// write back response data.
	response.SuccessWithData(filePart)
}

// check file chunks, mainly used to judge file integrity.
func checkChunkComplete(filePart *request.FilePartInfo) bool {
	currentChunkName := filePart.GetChunkFilename(filePart.CurrentCheckChunkNumber)
	exists := ioutil2.FileExists(currentChunkName)
	if exists {
		filePart.CurrentCheckChunkNumber++
		if filePart.CurrentCheckChunkNumber > filePart.GetTotalChunk() {
			// complete all transfers.
			return true
		}

		return checkChunkComplete(filePart)
	}
	return false
}

// get the chunk number collection that has been uploaded.
func getUploadedChunkNumbers(filePart *request.FilePartInfo) (isEqu bool, numbers []uint) {
	totalChunk := filePart.GetTotalChunk()
	var currentChunkNumber uint = 1
	uploadedChunkNumbers := make([]uint, 0)
	for {
		currentChunkName := filePart.GetChunkFilename(currentChunkNumber)
		exists := ioutil2.FileExists(currentChunkName)
		if exists {
			uploadedChunkNumbers = append(uploadedChunkNumbers, currentChunkNumber)
		}
		currentChunkNumber++
		if currentChunkNumber > totalChunk {
			break
		}
	}
	return len(uploadedChunkNumbers) == int(totalChunk), uploadedChunkNumbers
}
