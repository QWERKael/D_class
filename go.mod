module atk_D_class

go 1.13

replace (
	github.com/robfig/cron/v3 v3.0.1 => gopkg.in/robfig/cron.v3 v3.0.1
	gopkg.in/robfig/cron.v3 v3.0.1 => github.com/robfig/cron/v3 v3.0.1
)

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/c-bata/go-prompt v0.2.3
	github.com/go-ldap/ldap/v3 v3.2.4
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/term v0.0.0-20200520122047-c3ffed290a03 // indirect
	github.com/shirou/gopsutil v2.20.5+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	google.golang.org/grpc v1.30.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/robfig/cron.v3 v3.0.1
	gopkg.in/yaml.v2 v2.3.0
)
