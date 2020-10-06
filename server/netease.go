package server

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/ronething/music-push/config"
)

const netEaseBase = "https://music.163.com/#/song?id="

type NetEaseRank struct {
	Pre string
}

type NetEaseTracks struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Ar   []struct {
		ID    int           `json:"id"`
		Name  string        `json:"name"`
		Tns   []interface{} `json:"tns"`
		Alias []interface{} `json:"alias"`
	} `json:"ar"`
	Al struct {
		ID     int           `json:"id"`
		Name   string        `json:"name"`
		PicURL string        `json:"picUrl"`
		Tns    []interface{} `json:"tns"`
		PicStr string        `json:"pic_str"`
		Pic    int64         `json:"pic"`
	} `json:"al"`
	Dt          int   `json:"dt"`
	PublishTime int64 `json:"publishTime"`
}

type NetEasePlayList struct {
	Tracks                []NetEaseTracks `json:"tracks"`
	TrackNumberUpdateTime int64           `json:"trackNumberUpdateTime"`
	SubscribedCount       int             `json:"subscribedCount"`
	TrackCount            int             `json:"trackCount"`
	CommentThreadID       string          `json:"commentThreadId"`
	TrackUpdateTime       int64           `json:"trackUpdateTime"`
	PlayCount             int64           `json:"playCount"`
	Description           string          `json:"description"`
	Name                  string          `json:"name"`
	ID                    int             `json:"id"`
	ShareCount            int             `json:"shareCount"`
	CoverImgIDStr         string          `json:"coverImgId_str"`
	CommentCount          int             `json:"commentCount"`
}

type NetEasePlayListDetail struct {
	Code     int             `json:"code"`
	Playlist NetEasePlayList `json:"playlist"`
}

func (n *NetEaseRank) getNetEaseHost() string {
	return config.Config.GetString("netease.host")
}

// 云音乐飙升榜前 10
func (n *NetEaseRank) GetTop10() (string, error) {
	resp, err := req.Get(fmt.Sprintf("%s/playlist/detail?id=19723756&s=0", n.getNetEaseHost()))
	if err != nil {
		return "", err
	}
	var p NetEasePlayListDetail
	if err = resp.ToJSON(&p); err != nil {
		return "", err
	}

	var length int
	var tracks []NetEaseTracks
	updateTime := time.Unix(p.Playlist.TrackUpdateTime/1000, 0)
	update := updateTime.Format("2006-01-02")
	now := time.Now().Format("2006-01-02")
	log.Printf("update is %v, now is %v, n.Pre is %v\n,", update, now, n.Pre)
	if now != update || n.Pre == now { // 要么是不是当天，要么就已经是更新过了的
		return "", errors.New("get top 10 err")
	}
	if len(p.Playlist.Tracks) >= 7 {
		length = 7
	} else if len(p.Playlist.Tracks) <= 0 {
		return "", errors.New("tracks empty")
	} else {
		length = len(p.Playlist.Tracks)
	}
	for i := 0; i < length; i++ {
		tracks = append(tracks, p.Playlist.Tracks[i])
	}
	var content bytes.Buffer
	content.WriteString(fmt.Sprintf("云音乐飙升榜 (%s)\n\n", now))
	for i := 0; i < len(tracks); i++ {
		var artist []string
		ar := tracks[i].Ar
		if len(ar) > 0 {
			for j := 0; j < len(ar); j++ {
				artist = append(artist, ar[j].Name)
			}
		} else {
			artist = append(artist, "none")
		}
		content.WriteString(fmt.Sprintf("%d. %s - %s %s\n",
			i+1,
			tracks[i].Name,
			strings.Join(artist, "/"),
			fmt.Sprintf("%s%d", netEaseBase, tracks[i].ID),
		))
	}

	return content.String(), nil
}
