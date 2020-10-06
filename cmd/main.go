package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ronething/music-push/config"
	"github.com/ronething/music-push/pkg/wechat"
	"github.com/ronething/music-push/server"
)

var (
	filePath string // 配置文件路径
	help     bool   // 帮助
)

func usage() {
	fmt.Fprintf(os.Stdout, `music-push - music rank push service
Usage: music [-h help] [-c ./config.yaml]
Options:
`)
	flag.PrintDefaults()
}

func main() {

	flag.StringVar(&filePath, "c", "./config.yaml", "配置文件所在")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.Usage = usage
	flag.Parse()
	if help {
		flag.PrintDefaults()
		return
	}

	// 设置配置文件和静态变量
	config.SetConfig(filePath)

	// 登录微信
	err, loginMap := wechat.WechatLogin()
	if err != nil {
		log.Printf("登录微信发生错误, err: %v", err)
		return
	}
	users := wechat.GetSendUsers(loginMap, config.Config.GetStringSlice("cron.push"))

	scheduler := server.NewScheduler()
	scheduler.InitJob(loginMap, users)
	scheduler.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// DONE: 优雅关停
	for {
		s := <-c
		log.Printf("[main] 捕获信号 %s", s.String())
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// 停止调度器 并等待正在 running 的任务执行结束 TODO: 有没有必要设置一个 timeout 假设一直不停止怎么办
			ctx := scheduler.Stop()
			timer := time.NewTimer(1 * time.Second)
			for {
				select {
				case s = <-c: // 再次接收到中断信号 则直接退出
					if s == syscall.SIGINT {
						log.Printf("[main] 再次接收到退出信号 %s", s.String())
						goto End
					}
				case <-ctx.Done():
					log.Printf("[main] 调度器所有任务执行完成")
					goto End
				case <-timer.C:
					log.Printf("[main] 调度器有任务正在执行中")
					timer.Reset(1 * time.Second)
				}
			}
		End:
			timer.Stop()
			return // 很重要 不然程序无法退出
		case syscall.SIGHUP:
			log.Printf("[main] 终端断开信号，忽略")
		default:
			log.Printf("[main] other signal")
		}
	}

}
