package common

import (
	"fmt"
	"log"
	"os"
)


var Outfile *os.File

// 这是定义的日志的等级级别可根据自己的需求，定义自己需要的级别
var (
	Error *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
)

func init() {
	var err error
	Outfile, err = os.OpenFile("./log/agent.log",  os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	//fmt.Println(reflect.TypeOf(outfile))
	if err != nil {
		log.Panicf("open log file fail:%s ", err)
	}
	Error = log.New(Outfile, "[error]", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(Outfile, "[info]", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(Outfile, "[warn]", log.Ldate|log.Ltime|log.Lshortfile)
	//Error = log.New(Outfile, "[error]", log.Ldate|log.Ltime)
	//Info = log.New(Outfile, "[info]", log.Ldate|log.Ltime)
	//Warn = log.New(Outfile, "[warn]", log.Ldate|log.Ltime)
}

func CheckErr(funcName string, err error) {
	if err != nil {
		fmt.Println(funcName, err)
		panic(err)
	}
}

