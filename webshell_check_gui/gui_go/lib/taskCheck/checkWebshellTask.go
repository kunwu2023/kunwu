package taskCheck

import (
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/scan"
)

func checkWebshellTask(dirPath []remoteScanPath, taskBaseID int64, cloudType bool) error {
	db.Model(&dbpos.TaskBase{}).Where("id = ?", taskBaseID).Update("status", 2)
	for _, pathMap := range dirPath {
		isDirFlag, _ := scan.IsDir(pathMap.Path)
		taskType := scan.TaskType{
			IsQuickScan: false,
			TaskBaseID:  taskBaseID,
			IsCloud:     cloudType,
		}
		if isDirFlag {
			scan.BulkScanDir(pathMap.Path, taskType)
		} else {
			scan.BulkScanFile([]string{pathMap.Path}, taskType)
		}
	}
	db.Model(&dbpos.TaskBase{}).Where("id = ?", taskBaseID).Update("status", 3)
	return nil
}
