// 实现ORM框架的日志器，控制日志输出。

package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (

	errorLog = log.New(os.Stdout,"\033[31m[error]\033[0m ",log.LstdFlags|log.Lshortfile)  //os.Stdout是标准输出流函数，log.New创建一个新的Logger,红色输出
	infoLog = log.New(os.Stdout,"\033[34m[info ]\033[0m ",log.LstdFlags)  //蓝色输出，log.Lshortfile 支持显示文件名和代码行号
	loggers = []*log.Logger{errorLog,infoLog}
	mu sync.Mutex  //互斥锁

)

// log 方法:两种颜色，4种方法
var (
	Error = errorLog.Println  
	Errorf = errorLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

// log level 设置log日志层级：InfoLevel,ErrorLevel,Disabled

const (
	InfoLevel = iota  //1
	ErrorLevel   //2
	Disabled   //3   暂未使用
)

// 设置级别控制日志级别(日志分级，按照上面的const定义，分别是1~3级)

func SetLevel(level int) {
	mu.Lock() // 加个互斥锁
	defer mu.Lock()  //延时关闭

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)  //
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)  //输出会被定向到 ioutil.Discard，即不打印该日志
	}
}