module pcdn-server

go 1.24.3

require (
	github.com/bytedance/sonic v1.13.3
	github.com/go-playground/validator/v10 v10.26.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/liuhengloveyou/go-common v0.0.0-20250319112824-c28f82e5a12b
	github.com/liuhengloveyou/go-errors v0.0.0-20211025085721-e6717f3d23d1
	github.com/liuhengloveyou/passport v1.1.0
	github.com/qiniu/go-sdk/v7 v7.25.3
	github.com/redis/go-redis/v9 v9.10.0
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.36.6
	gorm.io/gorm v1.26.1
)

require (
	github.com/Blank-Xu/sqlx-adapter v1.0.1 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/robfig/cron/v3 v3.0.0 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1182 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.1115 // indirect
	golang.org/x/sync v0.15.0 // indirect
	gopkg.in/guregu/null.v4 v4.0.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/bmatcuk/doublestar/v4 v4.8.1 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/casbin/casbin/v2 v2.107.0 // indirect
	github.com/casbin/govaluate v1.7.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/google/uuid v1.6.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/liuhengloveyou/pcdn/protos v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.18.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.6.0
	xorm.io/builder v0.3.13 // indirect
)

replace github.com/liuhengloveyou/passport => ../../passport/

replace github.com/liuhengloveyou/pcdn/protos => ../protos/
