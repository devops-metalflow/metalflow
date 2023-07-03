package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"gorm.io/datatypes"
	"metalflow/pkg/global"
	"metalflow/pkg/grpc/securepb"
)

// SendSecureRequest send request to get security info from metalsecure.
func SendSecureRequest(address string, port int, reqBody *securepb.ServerRequest) (*securepb.ServerReply_Output, error) {
	conn, err := ConnectGrpc(address, port, context.Background())
	if err != nil {
		return nil, err
	}
	global.Log.Infof("grpc连接%s:%d成功", address, port)
	defer func(conn *grpc.ClientConn) {
		e := conn.Close()
		if e != nil {
			fmt.Println("grpc关闭服务失败")
		}
	}(conn)

	c := securepb.NewServerProtoClient(conn)
	r, err := c.SendServer(context.Background(), reqBody)
	if err != nil {
		global.Log.Errorf("[%s]metalsecure获取数据失败：%v", address, err)
		fmt.Printf("不能获取服务器%s的secure响应信息：%v", address, err)
		return nil, err
	}
	if outputErr := r.GetError(); outputErr != "" {
		global.Log.Errorf("[%s]metalsecure获取数据报错：%v", address, outputErr)
		fmt.Printf("获取服务器%s的secure信息报错：%s", address, outputErr)
		return nil, fmt.Errorf("err: %s", outputErr)
	}
	return r.GetOutput(), nil
}

const (
	apiVersion = "1.0.0"
	kind       = "metalsecure"
)

// GetImages get image stats from node
func GetImages(address string, port int) (datatypes.JSON, error) {
	imageReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Stats: &securepb.ServerRequest_Spec_Stats{Dockers: true},
		},
	}
	statsOutPut, err := SendSecureRequest(address, port, imageReq)
	if err != nil {
		return nil, err
	}
	stats, err := json.Marshal(statsOutPut.Stats.Dockers)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func GetDockerVul(address string, port int, dockers []*securepb.ServerRequest_Spec_Docker_Image) (datatypes.JSON, error) {
	serverReqSpecDocker := &securepb.ServerRequest_Spec_Docker{
		Vul: true,
	}
	if len(dockers) > 0 {
		serverReqSpecDocker = &securepb.ServerRequest_Spec_Docker{
			Images: dockers,
			Vul:    true,
		}
	}
	dockerVulReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Docker: serverReqSpecDocker,
		},
	}
	dockerVulOutPut, err := SendSecureRequest(address, port, dockerVulReq)
	if err != nil {
		global.Log.Errorf("[%s]获取metalsecure信息失败%v", address, err)
		return nil, err
	}
	jsonMap := make(map[string]any)
	err = json.Unmarshal([]byte(dockerVulOutPut.Docker.Vul), &jsonMap)
	if err != nil {
		global.Log.Errorf("[%s]解析metalsecure信息失败%v", address, err)
		return nil, err
	}
	return json.Marshal(jsonMap)
}

func GetDockerSBOM(address string, port int, dockers []*securepb.ServerRequest_Spec_Docker_Image) (datatypes.JSON, error) {
	serverReqSpecDocker := &securepb.ServerRequest_Spec_Docker{
		Sbom: true,
	}
	if len(dockers) > 0 {
		serverReqSpecDocker = &securepb.ServerRequest_Spec_Docker{
			Images: dockers,
			Sbom:   true,
		}
	}
	dockerSBOMReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Docker: serverReqSpecDocker,
		},
	}
	dockerSBOMOutPut, err := SendSecureRequest(address, port, dockerSBOMReq)
	if err != nil {
		return nil, err
	}
	jsonMap := make(map[string]any)
	err = json.Unmarshal([]byte(dockerSBOMOutPut.Docker.Sbom), &jsonMap)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonMap)
}

func GetBareVul(address string, port int, paths []string) (datatypes.JSON, error) {
	serverReqSpecBare := &securepb.ServerRequest_Spec_Bare{
		Vul: true,
	}
	if len(paths) > 0 {
		serverReqSpecBare = &securepb.ServerRequest_Spec_Bare{
			Paths: paths,
			Vul:   true,
		}
	}
	bareVulReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Bare: serverReqSpecBare,
		},
	}
	bareVulOutPut, err := SendSecureRequest(address, port, bareVulReq)
	if err != nil {
		return nil, err
	}
	jsonMap := make(map[string]any)
	err = json.Unmarshal([]byte(bareVulOutPut.Bare.Vul), &jsonMap)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonMap)
}

func GetBareSBOM(address string, port int, paths []string) (datatypes.JSON, error) {
	serverReqSpecBare := &securepb.ServerRequest_Spec_Bare{
		Sbom: true,
	}
	if len(paths) > 0 {
		serverReqSpecBare = &securepb.ServerRequest_Spec_Bare{
			Paths: paths,
			Sbom:  true,
		}
	}
	bareSBOMReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Bare: serverReqSpecBare,
		},
	}
	bareSBOMOutPut, err := SendSecureRequest(address, port, bareSBOMReq)
	if err != nil {
		return nil, err
	}
	jsonMap := make(map[string]any)
	err = json.Unmarshal([]byte(bareSBOMOutPut.Bare.Sbom), &jsonMap)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonMap)
}

func getAllSecureInfo(address string, port int) (dockerSbom, dockerVul, bareSbom, bareVul string, err error) {
	secureReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Bare: &securepb.ServerRequest_Spec_Bare{
				Sbom: true,
				Vul:  true,
			},
			Docker: &securepb.ServerRequest_Spec_Docker{
				Sbom: true,
				Vul:  true,
			},
		},
	}
	secureOutput, err := SendSecureRequest(address, port, secureReq)
	if err != nil {
		return
	}
	return secureOutput.Docker.Sbom, secureOutput.Docker.Vul, secureOutput.Bare.Sbom, secureOutput.Bare.Vul, nil
}

func GetRiskCount(address string, port int) (uint, error) {
	secureReq := &securepb.ServerRequest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Spec: &securepb.ServerRequest_Spec{
			Bare: &securepb.ServerRequest_Spec_Bare{
				Vul: true,
			},
			Docker: &securepb.ServerRequest_Spec_Docker{
				Vul: true,
			},
		},
	}
	secureOutput, err := SendSecureRequest(address, port, secureReq)
	if err != nil {
		return 0, err
	}
	var dockerRiskCount, bareRiskCount uint
	dockerRiskCount, _ = countSeverityRiskCount(secureOutput.Bare.Vul)
	bareRiskCount, _ = countSeverityRiskCount(secureOutput.Docker.Vul)
	return dockerRiskCount + bareRiskCount, nil
}

func GetSecureScore(address string, port int) (uint, error) {
	// 默认分数为满分
	const initScore = 100
	dockerSbom, dockerVul, bareSbom, bareVul, err := getAllSecureInfo(address, port)
	if err != nil {
		return 0, err
	}
	// 根据公式计算安全分数
	var (
		dockerSbomScore int
		dockerVulScore  int
		bareSbomScore   int
		bareVulScore    int
	)
	if dockerSbom != "" {
		dockerSbomScore, err = countLicense(dockerSbom)
		if err != nil {
			return 0, err
		}
	}
	if bareSbom != "" {
		bareSbomScore, err = countLicense(bareSbom)
		if err != nil {
			return 0, err
		}
	}
	if dockerVul != "" {
		dockerVulScore, err = countSeverity(dockerVul)
		if err != nil {
			return 0, err
		}
	}
	if bareVul != "" {
		bareVulScore, err = countSeverity(bareVul)
		if err != nil {
			return 0, err
		}
	}

	// 计算分数
	score := uint(initScore - (float64(bareSbomScore+dockerSbomScore)*0.25 + float64(bareVulScore+dockerVulScore)*0.75))
	return score, nil
}

func GetAllSeverityIDS(address string, port int) ([]string, error) {
	_, dockerVul, _, bareVul, err := getAllSecureInfo(address, port)
	if err != nil {
		return nil, err
	}
	var dockerSeverityIDS, bareSeverityIDS []string
	dockerSeverityIDS, err = getSeverityIDS(dockerVul)
	if err != nil {
		return nil, err
	}
	bareSeverityIDS, err = getSeverityIDS(bareVul)
	if err != nil {
		return nil, err
	}
	dockerSeverityIDS = append(dockerSeverityIDS, bareSeverityIDS...)
	return dockerSeverityIDS, nil
}

func countLicense(sbom string) (int, error) {
	sbomVals := make(map[string]*SbomItem)
	err := json.Unmarshal([]byte(sbom), &sbomVals)
	if err != nil {
		return 0, err
	}

	const weight = 50
	validLicenses := []string{"AGPL", "Basic Proprietary", "Commons Clause", "CPAL", "Elastic", "EUPL", "Oracle Java", "OSL"}

	// 遍历得到所有的licenses
	for _, v := range sbomVals {
		// 对各种许可证进行计数
		licenseCount := make(map[string]int)
		for _, artifact := range v.Artifacts {
			for _, license := range artifact.Licenses {
				licenseCount[license] += 1
			}
		}
		// 只计算有用的license
		var rate int
		for _, validLicense := range validLicenses {
			rate += licenseCount[validLicense]
		}
		// 只要计算次数不等于0，则满足要求直接返回
		if rate != 0 {
			return weight, nil
		}
	}
	return 0, nil
}

func countSeverityRiskCount(vul string) (uint, error) {
	vulVals := make(map[string]*VulItem)
	err := json.Unmarshal([]byte(vul), &vulVals)
	if err != nil {
		return 0, err
	}

	validSeverities := []string{"Critical", "High"}

	var count uint
	// 遍历得到所有的severity
	for _, v := range vulVals {
		// 对各种severity进行计数
		severityCount := make(map[string]uint)
		for _, match := range v.Matches {
			severityCount[match.Severity] += 1
		}
		// 只计算有用的severity
		for _, validSeverity := range validSeverities {
			count += severityCount[validSeverity]
		}
	}
	return count, nil
}

func getSeverityIDS(vul string) ([]string, error) {
	const (
		critical = "Critical"
		high     = "High"
	)

	vulVals := make(map[string]*VulItem)
	err := json.Unmarshal([]byte(vul), &vulVals)
	if err != nil {
		return nil, err
	}
	severityIDS := make([]string, 0)
	// 遍历得到所有的severity
	for _, v := range vulVals {
		// 得到severityID
		for _, match := range v.Matches {
			if match.Severity == critical || match.Severity == high {
				severityIDS = append(severityIDS, match.VulnerabilityID)
			}
		}
	}
	return severityIDS, nil
}

func countSeverity(vul string) (int, error) {
	vulVals := make(map[string]*VulItem)
	err := json.Unmarshal([]byte(vul), &vulVals)
	if err != nil {
		return 0, err
	}

	const weight = 50
	validSeverities := []string{"Medium", "High", "Critical"}

	// 遍历得到所有的severity
	for _, v := range vulVals {
		// 对各种severity进行计数
		severityCount := make(map[string]int)
		for _, match := range v.Matches {
			severityCount[match.Severity] += 1
		}
		// 只计算有用的severity
		var rate int
		for _, validSeverity := range validSeverities {
			rate += severityCount[validSeverity]
		}
		// 只要有一个镜像的计算总分数不等于0，则满足要求直接返回
		if rate != 0 {
			return weight, nil
		}
	}
	return 0, nil
}

type VulContent struct {
	Severity        string `json:"severity"`
	VulnerabilityID string `json:"vulnerabilityID"`
}

type VulItem struct {
	Matches []*VulContent `json:"matches"`
}

type SbomContent struct {
	Licenses []string `json:"licenses"`
}

type SbomItem struct {
	Artifacts []*SbomContent `json:"artifacts"`
}
