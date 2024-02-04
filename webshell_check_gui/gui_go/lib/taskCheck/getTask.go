package taskCheck

import "webshell_gui_go/lib/dbpos"

func GetTaskBase() []dbpos.TaskBase {
	var taskBase []dbpos.TaskBase
	db.Model(&dbpos.TaskBase{}).Find(&taskBase)
	return taskBase
}

func GetTaskList(id int64) []dbpos.TaskCheckList {
	var taskList []dbpos.TaskCheckList
	db.Model(&dbpos.TaskCheckList{}).Where("task_base_id = ? and (results != '正常' or cloud_results_flag = 4)", id).Find(&taskList)
	return taskList
}
