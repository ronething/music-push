module github.com/ronething/music-push

go 1.14

require (
	github.com/Han-Ya-Jun/qrcode2console v0.0.0-20190430081741-6890f5f0fdf5
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/imroc/req v0.3.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/spf13/viper v1.7.1
	itchat4go v0.0.0-00010101000000-000000000000
)

replace itchat4go => ./pkg/itchat4go
