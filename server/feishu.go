package server

import (
	"fmt"

	"github.com/imroc/req"
)

type FeiShu struct {
	WebHookUrl string
}

func NewFeiShu(webHookUrl string) *FeiShu {
	return &FeiShu{WebHookUrl: webHookUrl}
}

type FeiShuSendText struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func (f *FeiShu) Send(content string) error {
	return f.SendText(content)
}

func (f *FeiShu) SendText(content string) error {
	resp, err := req.Post(f.WebHookUrl, req.BodyJSON(&FeiShuSendText{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: content,
		},
	}))
	if err != nil {
		return err
	}
	fmt.Printf("resp is %v\n", resp.String())
	return err
}
