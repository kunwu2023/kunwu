package taskCheck

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/log"
	"webshell_gui_go/lib/scan"
)

var db *gorm.DB

func init() {
	db = dbpos.Db
}

type remoteScanPath struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	LastModifed int64  `json:"lastModified"`
}

type remoteScanConfig struct {
	ScanPath      []remoteScanPath `json:"scanPath"`
	PassWord      string           `json:"passWord"`
	CloudScanPath string           `json:"cloudScanPath"`
	ServerIp      string           `json:"serverIp"`
	SshPort       int              `json:"sshPort"`
	UserName      string           `json:"userName"`
	Model         int              `json:"model"`
	CloudScan     bool             `json:"cloud_scan"`
	DetectionMode bool             `json:"detection_mode"`
	MissionName   string           `json:"mission_name"`
}

type cloudRemoteScanConfig struct {
	ScanPath      []remoteScanPath `json:"scanPath"`
	PassWord      string           `json:"passWord"`
	CloudScanPath string           `json:"cloudScanPath"`
	ServerIp      string           `json:"serverIp"`
	SshPort       int              `json:"sshPort"`
	UserName      string           `json:"userName"`
	Model         int              `json:"model"`
	CloudScan     bool             `json:"cloud_scan"`
	DetectionMode bool             `json:"detection_mode"`
	MissionName   string           `json:"mission_name"`
}

func GetAllFilePaths(dirPath string) ([]string, error) {
	var paths []string
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			// 如果是目录，则递归获取所有文件路径
			subPaths, err := GetAllFilePaths(filePath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPaths...)
		} else {
			// 如果是文件，则将其路径添加到列表中
			paths = append(paths, filePath)
		}
	}
	return paths, nil
}

func AddCloudTask(request map[string]interface{}) (msg string, err error) {
	//fmt.Printf("Request: %v\n", request) // 打印请求内容
	jsonData := request["data"]

	// 将 jsonData 转换成 []byte 类型
	jsonDataStr, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}
	dataBytes := []byte(jsonDataStr)

	// 解析 JSON 数据到 remoteScanConfig 结构体中
	var config cloudRemoteScanConfig
	err = json.Unmarshal(dataBytes, &config)
	if err != nil {
		return "", err
	}
	if config.Model == 2 {
		// TODO 远程任务
		msg = "远程任务"
		taskBase := dbpos.TaskBase{
			Model:     int64(config.Model),
			DirPath:   "远程扫描：" + config.ServerIp,
			Status:    1,
			CreatedAt: time.Now().Unix(),
			Ip:        config.ServerIp,
			UserName:  config.UserName,
			PassWord:  config.PassWord,
			Port:      strconv.Itoa(config.SshPort),
		}
		db.Model(&dbpos.TaskBase{}).Create(&taskBase)
		SFTPClient := scan.NewSFTPClient(config.ServerIp, strconv.Itoa(config.SshPort), config.UserName, config.PassWord)
		db.Model(&taskBase).Update("status", 2)
		err = SFTPClient.ScanSshFiles(config.CloudScanPath, int(taskBase.ID))
		if err != nil {
			db.Model(&taskBase).Update("status", 4)
			db.Delete(&taskBase) // 删除创建的taskBase记录
			log.Logger.Error(fmt.Sprintf("ssh扫描错误：%s", err.Error()))
			return "", err
		} else {
			db.Model(&taskBase).Update("status", 3)
		}
	}
	return
}

func AddTask(request map[string]interface{}) (msg string, err error) {
	//fmt.Printf("Request: %v\n", request) // 打印请求内容

	cloudFlag, ok := request["cloudFlag"].(float64)
	if !ok {
		return "", fmt.Errorf("cloudFlag 类型断言失败")
	}
	var cloudType bool
	if cloudFlag == 1 {
		cloudType = true
	} else {
		cloudType = false
	}

	// 解析 JSON 字符串为 map[string]interface{}
	dataString := request["data"].(string)
	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(dataString), &dataMap)
	if err != nil {
		return "", err
	}

	// 将 dataMap 转换成 []byte 类型
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return "", err
	}

	// 解析 JSON 数据到 remoteScanConfig 结构体中
	var config remoteScanConfig
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return "", err
	}

	if config.Model == 1 {
		// TODO 本地任务
		msg = "本地任务"
		taskBase := dbpos.TaskBase{
			Model:     int64(config.Model),
			DirPath:   config.MissionName,
			Status:    1,
			CreatedAt: time.Now().Unix(),
		}
		db.Model(&dbpos.TaskBase{}).Create(&taskBase)
		//err = checkWebshellTask(config.ScanPath, taskBase.ID)
		//if err != nil {
		//	db.Model(&dbpos.TaskBase{}).Where("id = ?", taskBase.ID).Update("status", 3)
		//}
		go func() {
			err := checkWebshellTask(config.ScanPath, taskBase.ID, cloudType)
			if err != nil {
				db.Model(&dbpos.TaskBase{}).Where("id = ?", taskBase.ID).Update("status", 3)
			}
		}()
		idStr := strconv.FormatInt(taskBase.ID, 10)
		msg = "ID:" + idStr
	}
	return
}
