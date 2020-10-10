package server

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/ronething/music-push/config"
)

type Scheduler struct {
	C *cron.Cron
}

//NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	optLogs := cron.WithLogger(
		cron.VerbosePrintfLogger(
			log.New(os.Stdout, "[Cron]: ", log.LstdFlags)))

	c := cron.New(optLogs)
	return &Scheduler{C: c}

}

func (s *Scheduler) Run() {
	s.C.Start()
}

func (s *Scheduler) InitJob() {
	var err error
	n := NetEaseRank{ // pre init
		Pre: config.Config.GetString("webhook.pre"),
	}
	if _, err = s.C.AddFunc(config.Config.GetString("webhook.spec"), func() {
		now := time.Now().Format("2006-01-02")
		if n.Pre == now {
			log.Printf("当天已经推送过了")
			return
		}
		res, err := n.GetTop10()
		if err != nil {
			log.Printf("获取排行榜失败, err: %v\n", err)
			return
		}
		// webhook post
		feishuHooks := config.Config.GetStringSlice("webhook.feishu")
		success := 0 // 计数
		failed := 0
		for i := 0; i < len(feishuHooks); i++ {
			f := NewFeiShu(feishuHooks[i])
			if err = f.Send(res); err != nil {
				log.Printf("推送发生 err: %v\n", err)
				failed += 1
			} else {
				success += 1
			}
		}
		log.Printf("推送成功: %v, 推送失败: %v\n", success, failed)
		n.Pre = now // 设置标志
	}); err != nil {
		log.Printf("添加任务失败, err: %v\n", err)
		return
	}
}

func (s *Scheduler) Stop() context.Context {
	return s.C.Stop()
}
