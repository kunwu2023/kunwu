package taskCheck

import "webshell_gui_go/lib/dbpos"

func DelTask(data float64) {
	db.Model(&dbpos.TaskBase{}).Where("id = ?", data).Delete(&dbpos.TaskBase{})
	db.Model(&dbpos.TaskCheckList{}).Where("task_base_id = ?", data).Delete(&dbpos.TaskCheckList{})
}
