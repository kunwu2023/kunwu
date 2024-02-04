package scan

import (
	"io/ioutil"
	"os"
	"webshell_gui_go/lib/dbpos"
)

type ReadFileData struct {
	filePath string // filePath 文件路径
}

func NewReadFile(path string) *ReadFileData {
	return &ReadFileData{
		filePath: path,
	}
}

func (r *ReadFileData) ReadSshFile(taskID int) (fileData []byte, err error) {
	var taskBase dbpos.TaskBase
	db.Model(&dbpos.TaskBase{}).Where("id = ?", taskID).Find(&taskBase)

	client := NewSFTPClient(taskBase.Ip, taskBase.Port, taskBase.UserName, taskBase.PassWord)
	fileData, err = client.ReadFile(r.filePath)

	return
}

func (r ReadFileData) ReadFile() ([]byte, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	// 返回文件内容
	return content, nil
}

func (r ReadFileData) ReadFileToString() (string, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return "", err
	}

	// 将字节数组转换为字符串并返回
	return string(content), nil
}

func (r ReadFileData) DeleteFile() error {
	// 删除文件
	err := os.Remove(r.filePath)
	if err != nil {
		return err
	}
	return nil
}
