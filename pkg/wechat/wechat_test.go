package wechat

import (
	"fmt"
	s "itchat4go/service"
	"testing"
)

func TestWechatLogin(t *testing.T) {
	// login
	err, loginMap := WechatLogin()
	if err != nil {
		t.Fatalf("login wechat err: %v\n", err)
	}
	//WechatSendMsg("今天是个好日子，消息忽略", "")
	fmt.Println("开始获取联系人信息...")
	contactMap, err = s.GetAllContact(&loginMap)
	if err != nil {
		panic(err)
	}
	fmt.Printf("成功获取 %d 个 联系人信息,开始整理群组信息...\n", len(contactMap))

	for _, user := range contactMap {
		fmt.Printf("user is %v\n", user)
	}
}
