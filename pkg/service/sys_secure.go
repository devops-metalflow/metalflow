package service

import (
	"fmt"
	"gorm.io/datatypes"
	"io"
	"metalflow/models"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc"
	"metalflow/pkg/grpc/securepb"
	"metalflow/pkg/request"
	"strings"
)

func (s *MysqlService) GetRiskCountById(nodeId uint) (uint, error) {
	var metalSecure models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return 0, err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return 0, err
	}
	return grpc.GetRiskCount(node.Address, metalSecure.Port)
}

func (s *MysqlService) GetNodeImages(nodeId uint) (images datatypes.JSON, err error) {
	// get metalsecure info from database
	// TODO 因进行过worker初始化，所以id为2的worker就是metalsecure
	var metalSecure models.SysWorker
	err = s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return
	}
	// 直接连接grpc获取内容，不走异步任务
	images, err = grpc.GetImages(node.Address, metalSecure.Port)
	if err != nil {
		return
	}
	return images, nil
}

const (
	category                = "sbom"
	bareSecureFileName      = "bareSecure.sh"
	containerSecureFileName = "containerSecure.sh"
	fixRiskFileName         = "fixSecureRisk.sh"
)

func (s *MysqlService) GetNodeDockerSecureInfo(nodeId uint, req *request.SecureImage) (imagesReport datatypes.JSON, err error) {
	var metalSecure models.SysWorker
	err = s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return
	}
	// 直接进行grpc通信获取对应的报告内容
	images := make([]*securepb.ServerRequest_Spec_Docker_Image, 0)
	for _, image := range req.Images {
		images = append(images, &securepb.ServerRequest_Spec_Docker_Image{Repo: image.Repo, Tag: image.Tag})
	}
	if req.Category == category {
		imagesReport, err = grpc.GetDockerSBOM(node.Address, metalSecure.Port, images)
		if err != nil {
			return nil, err
		}
	} else {
		imagesReport, err = grpc.GetDockerVul(node.Address, metalSecure.Port, images)
		if err != nil {
			return nil, err
		}
	}
	return imagesReport, nil
}

func (s *MysqlService) GetNodeBareSecureInfo(nodeId uint, req *request.SecureBare) (bareReport datatypes.JSON, err error) {
	var metalSecure models.SysWorker
	err = s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return
	}

	if req.Category == category {
		bareReport, err = grpc.GetBareSBOM(node.Address, metalSecure.Port, req.Paths)
		if err != nil {
			return nil, err
		}
	} else {
		bareReport, err = grpc.GetBareVul(node.Address, metalSecure.Port, req.Paths)
		if err != nil {
			return nil, err
		}
	}
	return bareReport, nil
}

func (s *MysqlService) GetNodeSecureScore(nodeId uint) (score uint, err error) {
	var metalSecure models.SysWorker
	err = s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return
	}

	return grpc.GetSecureScore(node.Address, metalSecure.Port)
}

func (s *MysqlService) RunBareSecure(nodeId uint) error {
	ids := []uint{nodeId}
	secureShell := &secureShellInfo{
		content: "#!/bin/bash\necho 'hello world!'",
	}
	fileMetric := grpc.FileMetric{
		FilePath:   bareSecureFileName,
		RemoteDir:  remoteDir,
		IsRunnable: true,
		FileGetter: secureShell,
	}
	err := s.BatchUploadByIds(fileMetric, ids)
	if err != nil {
		global.Log.Errorf("执行裸金属安全修复失败：%v", err)
		return err
	}
	return nil
}

func (s *MysqlService) RunContainerSecure(nodeId uint) error {
	ids := []uint{nodeId}
	secureShell := &secureShellInfo{
		content: "#!/bin/bash\necho 'hello world!'",
	}
	fileMetric := grpc.FileMetric{
		FilePath:   containerSecureFileName,
		RemoteDir:  remoteDir,
		IsRunnable: true,
		FileGetter: secureShell,
	}
	err := s.BatchUploadByIds(fileMetric, ids)
	if err != nil {
		global.Log.Errorf("执行docker安全修复失败：%v", err)
		return err
	}
	return nil
}

func (s *MysqlService) FixSecureRisk(nodeId uint, cveId string) error {
	var metalSecure models.SysWorker
	err := s.TX.Model(new(models.SysWorker)).Where("id = ?", 3).First(&metalSecure).Error //nolint:gomnd
	if err != nil {
		return err
	}
	var node models.SysNode
	err = s.TX.Model(models.SysNode{}).Where("id = ?", nodeId).First(&node).Error
	if err != nil {
		return err
	}

	var severityIDS []string
	if cveId == "" {
		severityIDS, err = grpc.GetAllSeverityIDS(node.Address, metalSecure.Port)
		if err != nil {
			return err
		}
	} else {
		severityIDS = append(severityIDS, cveId)
	}
	// upload script file to node
	ids := []uint{nodeId}
	fixShellFile := &secureShellInfo{
		content: getFixShellContent(severityIDS),
	}
	fileMetric := grpc.FileMetric{
		FilePath:   fixRiskFileName,
		RemoteDir:  remoteDir,
		IsRunnable: true,
		FileGetter: fixShellFile,
	}
	err = s.BatchUploadByIds(fileMetric, ids)
	if err != nil {
		global.Log.Errorf("执行fix修复脚本失败：%v", err)
		return err
	}
	return nil
}

func getFixShellContent(ids []string) string {
	initShellContent := "#!/bin/bash\n"
	for _, id := range ids {
		shellUrl := fmt.Sprintf("https://factory/artifactory/example/"+
			"devops-metalflow/metalsecure/nvd/%s/%s.sh", id, id)
		downloadCmd := fmt.Sprintf("sudo curl -k "+
			"-uadmin:123456 -L %s -o /tmp/%s.sh\n",
			shellUrl, id)
		permissionCmd := fmt.Sprintf("sudo chmod +x /tmp/%s.sh\n", id)
		runCmd := fmt.Sprintf("sudo /tmp/%s.sh\n", id)
		initShellContent += downloadCmd
		initShellContent += permissionCmd
		initShellContent += runCmd
	}
	return initShellContent
}

type secureShellInfo struct {
	content string
}

func (c *secureShellInfo) GetFile() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(c.content)), nil
}
