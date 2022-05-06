package main

import (
	"cdp_agent/common"
	"cdp_agent/internal"
	"fmt"
	"time"
)

func main() {
	conf := common.Conf{}
	config := conf.GetModelClass("config.json")
	fmt.Printf("******************************Service启动:%s****************************\n", time.Now().Format("2006-01-02 15:04:05"))
	go internal.Service()
	fmt.Printf("******************************Agent启动:%s******************************\n", time.Now().Format("2006-01-02 15:04:05"))
	times := 1
	for {
		fmt.Printf("===========================进行第 %d 次数据采集，当前时间为:%s===========================\n", times, time.Now().Format("2006-01-02 15:04:05"))
		internal.Run() //数据采集
		time.Sleep(time.Second * time.Duration(int64(config.AgentCycle)))
		times += 1
	}
}
