package api

import (
	"fmt"
	"webshell_gui_go/lib/config"
	"webshell_gui_go/lib/dbpos"
	"webshell_gui_go/lib/log"
	"webshell_gui_go/lib/quickCheck"
	"webshell_gui_go/lib/scan"
	"webshell_gui_go/lib/taskCheck"
)

func MiddlewareProcessing(request map[string]interface{}) map[string]interface{} {
	// 处理请求并生成响应数据
	response := make(map[string]interface{})
	//log.Logger.Debug(fmt.Sprintf("request:%s", request))
	// 读取请求数据中的 `function` 字段，并根据其值执行相应的函数
	switch request["function"] {
	case "hello": // TODO test
		name, ok := request["data"].(string)
		if !ok {
			response["error"] = "invalid data"
			break
		}
		response["message"] = fmt.Sprintf("Hello, %s!", name)
		response["os"] = config.SystemType
		response["Arch"] = config.Arch
	case "quickCheckFile": // TODO 快速检测文件
		file, err := quickCheck.UpdateCheckFile(request)
		if err != nil {
			response["error"] = err
			log.Logger.Error(fmt.Sprintf("快速检测异常：%s", err.Error()))
		}
		response["file"] = file
		response["message"] = fmt.Sprintf("200_quickCheck")
	case "quickCheckList": // TODO 快速检测结果
		quickCheckList := quickCheck.GetList()
		response["quickCheckList"] = quickCheckList
		response["message"] = fmt.Sprintf("200_quickCheckList")
	case "addTask": // TODO 新建扫描任务
		msg, err := taskCheck.AddTask(request)
		if err != nil {
			response["error"] = err.Error()
			log.Logger.Error(fmt.Sprintf("创建任务异常：%s", err.Error()))
		}
		response["data"] = msg
		response["message"] = fmt.Sprintf("200_addTask")
	case "addCloudTask":
		msg, err := taskCheck.AddCloudTask(request)
		if err != nil {
			response["error"] = err.Error()
			log.Logger.Error(fmt.Sprintf("远程扫描异常：%s", err.Error()))
		}
		response["data"] = msg
		response["message"] = fmt.Sprintf("200_addCloudTask")
	case "getTask": // TODO 查询扫描任务列表
		taskBaseList := taskCheck.GetTaskBase()
		response["taskBaseList"] = taskBaseList
		response["message"] = fmt.Sprintf("200_getTask")
	case "getTaskList": // TODO 查询扫描任务列表
		data, _ := request["data"].(float64)
		intData := int64(data)
		taskBaseList := taskCheck.GetTaskList(intData)
		response["taskBaseList"] = taskBaseList
		response["message"] = fmt.Sprintf("200_getTaskList")
	case "getReadFile": // TODO 获取文件内容
		readFile := scan.NewReadFile(request["data"].(string))
		fileData, err := readFile.ReadFileToString()
		if err != nil {
			response["error"] = err.Error()
			log.Logger.Error(fmt.Sprintf("读取文件异常：%s", err.Error()))
		}
		response["fileData"] = fileData
		response["message"] = fmt.Sprintf("200_getReadFile")
	case "getSshReadFile": // TODO 获取ssh远程文件内容
		readFile := scan.NewReadFile(request["data"].(string))
		fileData, err := readFile.ReadSshFile(int(request["taskID"].(float64)))
		if err != nil {
			response["error"] = err.Error()
			log.Logger.Error(fmt.Sprintf("读取SSH文件异常：%s", err.Error()))
		}
		response["fileData"] = fileData
		response["message"] = fmt.Sprintf("200_getReadFile")
	case "delTask": // TODO 删除任务
		taskCheck.DelTask(request["data"].(float64))
	case "delFile": // TODO 删除文件
		readFile := scan.NewReadFile(request["data"].(string))
		err := readFile.DeleteFile()

		response["msg"] = "错误"
		if err != nil {
			response["error"] = err.Error()
			log.Logger.Error(fmt.Sprintf("删除文件异常：%s", err.Error()))
		} else {
			response["msg"] = "删除成功"
			// TODO 删除成功后修改数据库状态
			db := dbpos.Db
			db.Model(&dbpos.QuickCheckList{}).Where("path = ?", request["data"].(string)).Update("results", "已处置")
			db.Model(&dbpos.TaskCheckList{}).Where("path = ?", request["data"].(string)).Update("results", "已处置")
		}
		response["message"] = fmt.Sprintf("200_delFile")
	default:
		response["error"] = "invalid function"
	}
	//log.Logger.Debug(fmt.Sprintf("response:%s", response))
	return response
}
