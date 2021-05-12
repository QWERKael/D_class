BUILD_VERSION   := 1.0.0
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := go-version-sample
#SOURCE          := ./*.go
#TARGET_DIR      := /path-you-want/${BUILD_NAME}

all: cli serv show cmd mysql watcher sentry async echo transfer bee auth email config

cli:
	go build -ldflags \
	"-X atk_D_class/utils.BuildVersion=${BUILD_VERSION} \
	-X 'atk_D_class/utils.BuildTime=${BUILD_TIME}' \
	-X atk_D_class/utils.BuildName=cli" \
	-o target/cli client/client.go
serv:
	go build -ldflags \
	"-X atk_D_class/utils.BuildVersion=${BUILD_VERSION} \
	-X 'atk_D_class/utils.BuildTime=${BUILD_TIME}' \
	-X atk_D_class/utils.BuildName=serv" \
	-o target/serv server/server.go
show:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Show Plugin'" \
	-buildmode=plugin \
	-o target/plugin/show.so plugin/show.go
cmd:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=CMD Plugin'" \
	-buildmode=plugin \
	-o target/plugin/cmd.so plugin/cmd/cmd.go
mysql:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=MySQL Plugin'" \
	-buildmode=plugin \
	-o target/plugin/mysql.so plugin/mysql.go
watcher:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Watcher Plugin'" \
	-buildmode=plugin \
	-o target/plugin/watcher.so plugin/watcher/watcher.go
sentry:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Sentry Plugin'" \
	-buildmode=plugin \
	-o target/plugin/sentry.so plugin/watcher/sentry.go
async:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Async Plugin'" \
	-buildmode=plugin \
	-o target/plugin/async.so plugin/extra_task_manager/async.go
echo:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Echo Plugin'" \
	-buildmode=plugin \
	-o target/plugin/echo.so plugin/echo.go
transfer:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Transfer Plugin'" \
	-buildmode=plugin \
	-o target/plugin/transfer.so plugin/transfer/transfer.go
bee:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Bee Plugin'" \
	-buildmode=plugin \
	-o target/plugin/bee.so plugin/honeycomb/bee.go
auth:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Auth Plugin'" \
	-buildmode=plugin \
	-o target/plugin/auth.so plugin/auth/auth.go
email:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Auth Plugin'" \
	-buildmode=plugin \
	-o target/plugin/email.so plugin/email.go
config:
	go build -ldflags \
	"-X main.BuildVersion=${BUILD_VERSION} \
	-X 'main.BuildTime=${BUILD_TIME}' \
	-X 'main.BuildName=Auth Plugin'" \
	-buildmode=plugin \
	-o target/plugin/config.so plugin/config/config.go

clean:
	rm -f cli serv plugin/sentry.so plugin/watcher.so plugin/mysql.so plugin/cmd.so plugin/show.so plugin/bee.so plugin/auth.so plugin/email.so plugin/config.so

install:
    # mkdir -p ${TARGET_DIR}
    # cp ${BUILD_NAME} ${TARGET_DIR} -f

.PHONY : all cli serv show cmd mysql watcher sentry async echo transfer bee auth email config clean install ${BUILD_NAME}
