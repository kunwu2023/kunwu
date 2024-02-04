package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"runtime"
)

var Logger = InitLogger()
var logBash string

func initLog() (LogBash string, err error) {
	homeDir, err := os.UserHomeDir()
	// 在用户的 "Documents" 文件夹下创建一个文件夹存放数据库文件
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		LogBash = filepath.Join(homeDir, "Library", "kw_webshell_scan_DATA")
	} else {
		LogBash = filepath.Join(homeDir, "Documents", "kw_webshell_scan_DATA")
	}
	err = os.MkdirAll(LogBash, os.ModePerm)
	if err != nil {
		return
	}
	fmt.Println("日志文件路径：" + LogBash)
	return
}

func InitLogger() *zap.Logger {

	//获取编码器
	encoderConfig := zap.NewProductionEncoderConfig() //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	//encoderConfig := zap.NewDevelopmentEncoderConfig()                                //NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder                      //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	//encoderConfig.EncodeCaller = zapcore.FullCallerEncoder //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	if logBash == "" {
		logBash, _ = initLog()
	}
	fileName := logBash + "/kw.log"
	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName, //日志文件存放目录
		MaxSize:    50,       //文件大小限制,单位MB
		MaxBackups: 5,        //最大保留日志文件数量
		MaxAge:     30,       //日志文件保留天数
		Compress:   false,    //是否压缩处理
	})
	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
	fileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)

	return zap.New(fileCore, zap.AddCaller()) //AddCaller()为显示文件名和行号
}
