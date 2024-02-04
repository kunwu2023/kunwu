package quickCheck

import (
	"fmt"
	"gorm.io/gorm"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/scan"
)

var db *gorm.DB

func init() {
	db = dbpos.Db
}

type fileInfo struct {
	LastModified    int64  `json:"lastModified"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	Size            int64  `json:"size"`
	DirType         string `json:"dirType"`
	DetectionResult string `json:"detectionResult"`
}

func UpdateCheckFile(request map[string]interface{}) ([]fileInfo, error) {
	cloudFlag := request["cloudFlag"].(float64)
	var cloudType bool
	if cloudFlag == 1 {
		cloudType = true
	} else {
		cloudType = false
	}
	var files []fileInfo
	dataArr, ok := request["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data format")
	}
	for _, item := range dataArr {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		file := fileInfo{
			LastModified: int64(itemMap["lastModified"].(float64)),
			Name:         itemMap["name"].(string),
			Path:         itemMap["path"].(string),
			Size:         int64(itemMap["size"].(float64)),
		}
		file.DetectionResult = "待检测"
		files = append(files, file)
	}
	for _, v := range files {
		isDirFlag, _ := scan.IsDir(v.Path)
		taskType := scan.TaskType{
			IsQuickScan: true,
			TaskBaseID:  0,
			IsCloud:     cloudType,
		}
		if isDirFlag {
			scan.BulkScanDir(v.Path, taskType)
		} else {
			scan.BulkScanFile([]string{v.Path}, taskType)
		}
	}
	return files, nil
}
