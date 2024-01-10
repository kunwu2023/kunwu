package scan

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kunwu_cli/lib/config"
	"net/http"
	"os"
	"time"
)

var client *http.Client

func init() {
	// 忽略 https 证书验证
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: transport}
}

type Scanner struct {
	Path          string // path 扫描路径
	Model         int    // model 扫描模式：1文件模式，2文件夹模式
	Size          int64  // size 文件大小
	ModifyTime    int64  // modifyTime 最后一次修改时间
	CloudScanFlag bool   // CloudScanFlag 是否进行云端扫描，true表示进行云端扫描，false表示不进行云端扫描

}

func NewScanner(path string, cloudScanFlag bool) *Scanner {
	size, modifyTime, err := getFileSizeAndModifyTime(path)
	if err != nil {
		fmt.Println("文件读取异常", err.Error())
	}
	isDirFlag, err := IsDir(path)
	if err != nil {
		fmt.Println("文件读取异常", err.Error())
	}
	var model int
	if isDirFlag {
		model = 2
	} else {
		model = 1
	}
	return &Scanner{
		Path:          path,
		Model:         model,
		ModifyTime:    modifyTime,
		Size:          size,
		CloudScanFlag: cloudScanFlag,
	}
}

func (s *Scanner) CloudScan() (lastResults string, err error) {
	cloudLastResults, err := s.cloudScan()
	// 下面是云端检测
	if err != nil {
		//fmt.Println("云端扫描 ERROR", err.Error())
	} else {
		lastResults = cloudLastResults
	}
	return
}

func (s *Scanner) setErrorCloudScan() {
	//fmt.Println(s.Path, "云检查异常")
	return
}

func (s *Scanner) setCloudScan(results string) string {
	if results == "云端引擎超时" {
		strings := []string{"恶意", "正常"}
		return strings[1]
	}
	return ""
}

// cloudScan 云端扫描
func (s *Scanner) cloudScan() (cloudReslts string, err error) {
	content, err := ioutil.ReadFile(s.Path)
	if err != nil {
		fmt.Println("读取文件内容失败", err)
		s.setErrorCloudScan()
		return
	}
	fileMD5 := calculateFileMD5(string(content))
	// TODO 先查一遍云端有没有，有的话直接返回结果
	cloudResult, err := getCloudResult(fileMD5)
	if err != nil {
		s.setErrorCloudScan()
		return "", err
	}
	if cloudResult == "正常" || cloudResult == "恶意" {
		// TODO 有结果
		s.setCloudScan(cloudResult)
		return cloudResult, nil
	}
	type RequestBody struct {
		WebshellText string `json:"webshell_text"`
	}
	requestBody := RequestBody{WebshellText: base64.StdEncoding.EncodeToString(content)}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("构建请求体失败", err)
		s.setErrorCloudScan()
		return
	}
	resp, err := client.Post(config.Url+"/api/v1/anonymous_up_file?apikey="+config.ApiKey, "application/json", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		//fmt.Println("请求云端扫描接口失败", err)
		s.setErrorCloudScan()
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取云端扫描响应体失败", err)
		s.setErrorCloudScan()
		return
	}

	cloudResult, err = s.getCloudResults(fileMD5)
	if err != nil {
		s.setErrorCloudScan()
		return
	}
	s.setCloudScan(cloudResult)
	return cloudResult, nil
}

func (s *Scanner) getCloudResults(md5 string) (result string, err error) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 12; i++ {
		select {
		case <-ticker.C:
			result, err = getCloudResult(md5)
			if err != nil {
				return
			}
			if result == "正常" || result == "恶意" {
				return
			}
		}
	}

	// 如果请求了12次后 Result 值仍然为 0，则将 Result 值设置为 4 并输出
	result = "云端引擎超时"
	return
}

func getCloudResult(md5 string) (result string, err error) {
	type response struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result string `json:"result"`
	}
	url := config.Url + "/api/v1/anonymous_see_file?apikey=" + config.ApiKey + "&md5=" + md5
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respData response
	var jsonErr error
	if jsonErr = json.Unmarshal(body, &respData); jsonErr != nil {
		err = jsonErr
		return
	}

	result = respData.Result
	return
}

func calculateFileMD5(str string) (retMd5 string) {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func IsDir(path string) (bool, error) {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return dirInfo.IsDir(), nil
}

func getFileSizeAndModifyTime(path string) (int64, int64, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	// 获取文件大小
	size := fileInfo.Size()

	// 获取最后修改时间
	modifyTime := fileInfo.ModTime().Unix()

	return size, modifyTime, nil
}
