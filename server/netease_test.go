package server

import (
	"testing"

	"github.com/ronething/music-push/config"
)

func TestNetEaseRank_GetTop10(t *testing.T) {
	config.SetConfig("/Users/ronething/Documents/music-push/config/dev.yaml")
	n := NetEaseRank{}
	s, err := n.GetTop10()
	if err != nil {
		t.Errorf("err: %v\n", err)
		return
	}
	t.Logf("res is %v\n", s)
}
