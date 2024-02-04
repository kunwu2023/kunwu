package quickCheck

import "webshell_gui_go/lib/dbpos"

func GetList() []dbpos.QuickCheckList {
	var quickList []dbpos.QuickCheckList
	db.Model(&dbpos.QuickCheckList{}).Find(&quickList)
	return quickList
}
