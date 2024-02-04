package scan

import (
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/log"
)

type TaskType struct {
	IsQuickScan bool  // 是否是快速扫描
	TaskBaseID  int64 // 本地扫描的任务ID
	IsCloud     bool  // 是否开启云端扫描
}

func BulkScanFile(paths []string, taskType TaskType) {
	log.Logger.Info("Cloud Scan start", zap.Int("file count", len(paths)))
	BulkCloudScan(paths, taskType)
}

func BulkScanDir(path string, taskType TaskType) {
	f, err := os.Stat(path)
	if err != nil {
		log.Logger.Error("open file error: ", zap.Error(err))
	}
	var paths []string
	if f.IsDir() {
		paths, err = getAllFiles(path)
		if err != nil {
			log.Logger.Error("get all files error: ", zap.Error(err))
		}
	} else {
		paths = []string{path}
	}
	BulkScanFile(paths, taskType)
}

func getAllFiles(path string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() { // 判断是否为普通文件
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// CreateWebshellRecord 创建快速扫描的webshell记录
func CreateWebshellRecord(path string, taskType TaskType, cloudScan ...bool) {
	log.Logger.Debug("create webshell record", zap.String("path", path))
	f, err := os.Stat(path)
	if err != nil {
		log.Logger.Error("os Stat err", zap.Error(err))
		return
	}

	// 设置 cloudScan 的默认值
	cs := false
	if len(cloudScan) > 0 {
		cs = cloudScan[0]
	}
	var cloudResultsFlag int64 // 是否是云端扫描的结果
	if cs {
		cloudResultsFlag = 3
	} else {
		cloudResultsFlag = 1
	}

	if taskType.IsQuickScan {
		quickCheck := dbpos.QuickCheckList{
			Path:             path,
			Results:          "恶意",
			Size:             strconv.Itoa(int(f.Size())),
			CloudResultsFlag: cloudResultsFlag,
			ModificationTime: f.ModTime().Unix(),
		}
		db.Model(dbpos.QuickCheckList{}).Create(&quickCheck)
	} else {
		taskCheckList := dbpos.TaskCheckList{
			TaskBaseID:       taskType.TaskBaseID,
			Path:             path,
			Results:          "恶意",
			CloudResultsFlag: cloudResultsFlag,
			ModificationTime: f.ModTime().Unix(),
			Size:             strconv.Itoa(int(f.Size())),
		}
		db.Model(&dbpos.TaskCheckList{}).Create(&taskCheckList)
	}
}

// InArray 判断是否存在数组中
//func InArray(needle string, haystack []string) bool {
//	for _, v := range haystack {
//		if needle == v {
//			return true
//		}
//	}
//	return false
//}
