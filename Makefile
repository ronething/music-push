# 构建脚本

# set-env copy-config 在这里被依赖 在 build-master 和 build-worker 也被依赖，但是不会执行两次
.PHONY: deploy
deploy: set-env copy-config build upx

.PHONY: build
build: set-env copy-config
	go build -v -o bin/music cmd/main.go
	@echo "build music success"

.PHONY: copy-config
copy-config:
	rm -rf bin && mkdir -p bin && cp config/*.yaml bin/
	@echo "copy config success"

.PHONY: set-env
set-env:
	export GO111MODULE=on
	export GOPROXY=https://goproxy.io
	@echo "set env success"

# NOTICE: 需要确保有安装 upx
.PHONY: upx
upx:
	upx -v bin/music
