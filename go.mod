module fuse_file_system

go 1.14

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200403201458-baeed622b8d8

replace golang.org/x/tools => github.com/golang/tools v0.0.0-20200403190813-44a64ad78b9b

replace golang.org/x/text => github.com/golang/text v0.3.2

replace golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543

replace golang.org/x/exp => github.com/golang/exp v0.0.0-20200331195152-e8c3332aa8e5

replace golang.org/x/mod => github.com/golang/mod v0.2.0

replace golang.org/x/net => github.com/golang/net v0.0.0-20200324143707-d3edc9973b7e

replace golang.org/x/sync => github.com/golang/sync v0.0.0-20200317015054-43a5402ce75a

replace golang.org/x/sys => github.com/golang/sys v0.0.0-20200331124033-c3d80250170d

require (
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/hanwen/go-fuse/v2 v2.0.3
	github.com/ipfs/go-log/v2 v2.1.1
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.6.0
	github.com/tebeka/strftime v0.1.4 // indirect
)
