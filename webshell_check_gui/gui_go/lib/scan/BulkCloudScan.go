package scan

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"webshell_gui_go/lib/config"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/log"
)

type CloudResultType struct {
	Webshell []string `json:"webshell"`
	Normal   []string `json:"normal"`
}

func BulkCloudScan(path []string, taskType TaskType) {
	var fileMd5List []string
	var md5Map = make(map[string]string)  // [md5]path
	var pathMap = make(map[string]string) // [path]md5
	var cloudResult CloudResultType

	// 先查一遍云端有没有，有的话直接存入结果
	for _, v := range path {
		db.Model(&dbpos.QuickCheckList{}).Where("path = ?", v).Update("cloud_results_flag", 2)
		db.Model(&dbpos.TaskCheckList{}).Where("path = ?", v).Update("cloud_results_flag", 2)
		content, err := os.ReadFile(v)
		if err != nil {
			log.Logger.Error("读取文件内容失败", zap.String("path", v), zap.Error(err))
			continue
		}
		m := calculateFileMD5(content)
		md5Map[m] = v
		pathMap[v] = m
		fileMd5List = append(fileMd5List, m)
	}
	cloudResult, err := BulkGetCloudResult(fileMd5List)
	if err != nil {
		log.Logger.Error("获取云端扫描结果失败", zap.Error(err))
		return
	}
	// 云端有结果的文件
	log.Logger.Info("云端有结果的文件", zap.Int("webshell", len(cloudResult.Webshell)), zap.Int("normal", len(cloudResult.Normal)))
	for _, v := range cloudResult.Webshell {
		CreateWebshellRecord(md5Map[v], taskType, true)
	}

	//找出 云端没有返回结果的文件
	var CloudFileMap = make(map[string]string)
	var trackTable []string
	for _, v := range path {
		if !InArray(pathMap[v], cloudResult.Webshell) && !InArray(pathMap[v], cloudResult.Normal) {
			CloudFileMap[pathMap[v]] = v
			trackTable = append(trackTable, pathMap[v])
		}
	}
	log.Logger.Info("云端没有返回结果的文件", zap.Int("CloudFileMap", len(CloudFileMap)))
	if len(CloudFileMap) != 0 {
		composer, err := crateZipComposer(CloudFileMap)
		if err != nil {
			log.Logger.Error("crateZipComposre err", zap.Error(err))
			return
		}
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		// 创建一个表单文件
		fileWriter, err := writer.CreateFormFile("file", "webshell.zip")
		if err != nil {
			log.Logger.Error("writer.CreateFormFile err", zap.Error(err))
			return
		}
		// 写入文件数据
		_, err = io.Copy(fileWriter, bytes.NewReader(composer))
		if err != nil {
			log.Logger.Error("io.Copy err", zap.Error(err))
			return
		}
		writer.Close()
		req, err := http.NewRequest("POST", config.Url+"/api/v1/anonymous_up_zip_file?apikey="+config.ApiKey, buf)
		if err != nil {
			log.Logger.Error("http.NewRequest err", zap.Error(err))
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Logger.Error("http.DefaultClient.Do err", zap.Error(err))
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		log.Logger.Debug("body", zap.String("body", string(body)))

		go func(t []string, m map[string]string, number int) {
			var n int
			for true {
				result, err := BulkGetCloudResult(t)
				if err != nil {
					log.Logger.Error("BulkGetCloudResult err", zap.Error(err))
					return
				}
				if len(result.Webshell) != 0 || len(result.Normal) != 0 {
					for _, v := range result.Webshell {
						CreateWebshellRecord(m[v], taskType, true)
						t = removeElement(t, v)
					}
					for _, v := range result.Normal {
						t = removeElement(t, v)
					}
					n = 0
				} else {
					n++
					if n == number {
						log.Logger.Info("云端扫描超时", zap.Error(errors.New("云端扫描超时")))
						return
					}
				}
				if len(t) == 0 {
					log.Logger.Info("云端扫描结束", zap.Error(errors.New("云端扫描结束")))
					return
				}
				time.Sleep(5 * time.Second)
			}
		}(trackTable, md5Map, 6)

		time.Sleep(time.Second * 10)
	}
}

func BulkGetCloudResult(md5List []string) (result CloudResultType, err error) {
	type response struct {
		Code   int             `json:"code"`
		Msg    string          `json:"msg"`
		Result CloudResultType `json:"result"`
	}
	b, err := json.Marshal(md5List)
	if err != nil {
		log.Logger.Error("json.Marshal(md5List) err", zap.Error(err))
	}
	req, err := http.NewRequest("POST", config.Url+"/api/v2/anonymous_see_file?apikey="+config.ApiKey, bytes.NewBuffer(b))
	if err != nil {
		log.Logger.Error("http.NewRequest err", zap.Error(err))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Logger.Error("http.DefaultClient.Do err", zap.Error(err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error("http.DefaultClient.Do err", zap.Error(err))
		return
	}

	var respData response
	var jsonErr error
	if jsonErr = json.Unmarshal(body, &respData); jsonErr != nil {
		err = jsonErr
		return
	}
	if respData.Code == -1 {
		err = errors.New(respData.Msg)
		return
	}
	result = respData.Result
	return
}

// InArray 判断是否存在数组中
func InArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func crateZipComposer(md5Map map[string]string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	for k, v := range md5Map {
		f, err := zipWriter.Create(k)
		if err != nil {
			log.Logger.Error("zipWriter Create err", zap.Error(err))
			continue
		}
		content, err := os.ReadFile(v)
		if err != nil {
			log.Logger.Error("os.ReadFile err", zap.Error(err))
			continue
		}
		_, err = f.Write(content)
		if err != nil {
			log.Logger.Error("f.Write err", zap.Error(err))
			continue
		}
	}
	err := zipWriter.Close()
	if err != nil {
		log.Logger.Error("zipWriter Close err", zap.Error(err))
		return nil, err
	}
	return buf.Bytes(), nil
}

func removeElement(nums []string, val string) []string {
	var newNums []string
	for _, num := range nums {
		if num != val {
			newNums = append(newNums, num)
		}
	}
	return newNums
}
