module agent

go 1.24.3

require github.com/liuhengloveyou/go-selfupdate v0.0.0-20230714125711-e1c03e3d6ac7

require (
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

require (
	github.com/kr/binarydist v0.1.0 // indirect
	github.com/robfig/cron/v3 v3.0.0 // indirect
	go.uber.org/zap v1.27.0
)

replace github.com/liuhengloveyou/go-selfupdate => ../../dev/go-selfupdate/
