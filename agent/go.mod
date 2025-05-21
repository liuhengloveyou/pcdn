module agent

go 1.24.3

require (
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/liuhengloveyou/go-selfupdate v0.0.0-20230714125711-e1c03e3d6ac7
	github.com/liuhengloveyou/pcdn/protos v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.6.0 // indirect
	github.com/go-redis/redis/v8 v8.8.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jmoiron/sqlx v1.3.3 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/liuhengloveyou/go-common v0.0.0-20210507010927-c33a8f7ce0fb // indirect
	github.com/liuhengloveyou/go-errors v0.0.0-20210427064226-f3c31c856a7a // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.opentelemetry.io/otel v0.20.0 // indirect
	go.opentelemetry.io/otel/metric v0.20.0 // indirect
	go.opentelemetry.io/otel/trace v0.20.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/sys v0.0.0-20210503173754-0981d6026fa6 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/guregu/null.v3 v3.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/kr/binarydist v0.1.0 // indirect
	github.com/liuhengloveyou/passport v1.1.0
	github.com/robfig/cron/v3 v3.0.0
	go.uber.org/zap v1.27.0
)

replace (
	github.com/liuhengloveyou/go-selfupdate => ../../dev/go-selfupdate/
	github.com/liuhengloveyou/pcdn/protos => ../protos/
)
