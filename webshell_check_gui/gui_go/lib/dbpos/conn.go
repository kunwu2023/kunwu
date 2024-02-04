package dbpos

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"webshell_gui_go/lib/log"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	Db, err = initDB("webshell_db.sqlite3")
	if err != nil {
		log.Logger.Error(fmt.Sprintf("数据库起始化失败：%s", err.Error()))
		panic(err)
	}
}

func initDB(dbFile string) (*gorm.DB, error) {
	// 获取当前用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户主目录：%v", err)
	}

	// 在用户的 "Documents" 文件夹下创建一个文件夹存放数据库文件
	var appDataDir string
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		appDataDir = filepath.Join(homeDir, "Library", "kw_webshell_scan_DATA")
	} else {
		appDataDir = filepath.Join(homeDir, "Documents", "kw_webshell_scan_DATA")
	}
	err = os.MkdirAll(appDataDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("无法创建应用数据目录：%v", err)
	}

	// 将数据库文件路径设置为新创建的文件夹下
	dbFilePath := filepath.Join(appDataDir, dbFile)

	log.Logger.Debug(fmt.Sprintf("数据库文件路径：%s", dbFilePath))

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		// 如果不存在，创建数据库文件
		file, err := os.Create(dbFilePath)
		if err != nil {
			return nil, fmt.Errorf("创建数据库文件失败：%v", err)
		}
		file.Close()
	}

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移数据表
	err = db.AutoMigrate(&QuickCheckList{}, &TaskBase{}, &TaskCheckList{})
	if err != nil {
		return nil, fmt.Errorf("创建表失败：%v", err)
	}

	// 清空 quick_check_list 表
	if err := db.Exec("DELETE FROM quick_check_list").Error; err != nil {
		return nil, fmt.Errorf("清空 quick_check_list 表失败：%v", err)
	}

	return db, nil
}
