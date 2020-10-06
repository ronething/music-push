package server

import (
	"context"
	"fmt"
	m "itchat4go/model"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/ronething/music-push/config"
	"github.com/ronething/music-push/pkg/wechat"
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

func (s *Scheduler) InitJob(loginMap m.LoginMap, users []string, testUsers []string) {
	var err error
	_, err = s.C.AddFunc(config.Config.GetString("cron.spec"), func() {
		err = wechat.WechatSendMsgs(fmt.Sprintf("保活 %s", time.Now().String()), testUsers, loginMap)
		if err != nil {
			log.Printf("定时任务 - 发送微信消息发生错误 err: %v\n", err)
			return
		}
	})
	n := NetEaseRank{}
	_, err = s.C.AddFunc(config.Config.GetString("cron.spec"), func() {
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
		err = wechat.WechatSendMsgs(res, users, loginMap)
		if err != nil {
			log.Printf("云音乐排行榜任务 - 发送微信消息发生错误 err: %v\n", err)
			return
		}
		n.Pre = now // 设置标志
	})
}

func (s *Scheduler) Stop() context.Context {
	return s.C.Stop()
}
