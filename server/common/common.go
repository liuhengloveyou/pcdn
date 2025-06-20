package common

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm/logger"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	gocommon "github.com/liuhengloveyou/go-common"
	"github.com/liuhengloveyou/passport/accessctl"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	Slot           = "0x730x680x6f0x77"
	MAX_UPLOAD_LEN = 5242880
)

// redis key
const (
	AGENT_KEY_PREFIX         = "agent/"
	AGENT_MONITOR_KEY_PREFIX = "agent/monitor/"
	TASK_RESPONSE_KEY_PREFIX = "task/resp/"
	AGENT_TASK_KEY_PREFIX    = "agent/task/"
)

var confile = flag.String("c", "app.conf.yaml", "配置文件")
var (
	ProxyURL = "" // "http://127.0.0.1:10808" // 是否是开发环境

	ServConfig ConfigStruct

	Logger      *zap.Logger
	OrmCli      *gorm.DB
	RedisClient *redis.Client
)

type ConfigStruct struct {
	PID            string `yaml:"pid"`
	AccessName     string `yaml:"access_name"`
	Host           string `yaml:"host"`
	HttpServerAddr string `yaml:"http_server_addr"`
	TcpServerAddr  string `yaml:"tcp_server_addr"`
	PGURN          string `yaml:"pg_urn"`
	RedisAddr      string `yaml:"redis_addr"`
	UploadDir      string `yaml:"upload_dir"`
	FileDir        string `yaml:"file_dir"`
	LogDir         string `yaml:"log_dir"`
	LogLevel       string `yaml:"log_level"`
	ImgDir         string `yaml:"img_dir"`       // 图片上传目录
	SotreImgDir    string `yaml:"store_img_dir"` // 入库单图片上传目录
	// 收款地址
	Receiver string `yaml:"receiver"`

	// 管理员UID
	AdminUID int64 `yaml:"admin_id"`
}

func init() {
	if e := gocommon.LoadYamlConfig(*confile, &ServConfig); e != nil {
		panic(e)
	}
	if ServConfig.AccessName == "" {
		panic("access_name is empty")
	}

	if e := InitLog(ServConfig.LogDir, ServConfig.LogLevel); e != nil {
		panic(e)
	}

	if len(ServConfig.PGURN) > 0 {
		if e := InitGorm(ServConfig.PGURN); e != nil {
			fmt.Println("InitGorm: ", ServConfig.PGURN, e)
			panic(e)
		}
	}

	if len(ServConfig.RedisAddr) > 0 {
		if e := InitRedis(ServConfig.RedisAddr); e != nil {
			panic(e)
		}
	}
}

func InitLog(logDir, logLevel string) error {
	writer, _ := rotatelogs.New(
		logDir+"log.%Y%m%d%H%M",
		rotatelogs.WithLinkName("log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	level := zapcore.DebugLevel
	if e := level.UnmarshalText([]byte(logLevel)); e != nil {
		return e
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),
		zapcore.AddSync(writer),
		level)

	Logger = zap.New(core, zap.AddCaller())

	Logger.Info("initLog OK\n")
	return nil
}

func InitGorm(dsn string) (e error) {

	if OrmCli, e = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,  // "host=localhost user=pcdn password=pcdn12321 dbname=pcdn port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Warn),
	}); e != nil {
		return
	}

	sqlDB, err := OrmCli.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return
}

func InitRedis(redisAddr string) (e error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 测试连接
	e = RedisClient.Ping(context.Background()).Err()
	if e != nil {
		Logger.Sugar().Errorf("连接Redis失败: %v", e)
		RedisClient = nil
		return
	}

	return nil
}

func IsNowOpenningTime(timeArr []string) bool {
	if len(timeArr) != 2 {
		return true // 没有限制
	}

	a0 := strings.Split(timeArr[0], ":")
	a1 := strings.Split(timeArr[1], ":")
	if len(a0) != 2 || len(a1) != 2 {
		return true
	}

	t00, _ := strconv.Atoi(a0[0])
	t01, _ := strconv.Atoi(a0[1])
	t10, _ := strconv.Atoi(a1[0])
	t11, _ := strconv.Atoi(a1[1])
	t0 := t00*60 + t01
	t1 := t10*60 + t11
	tnow := time.Now().Hour()*60 + time.Now().Minute()

	if t1 > t0 {
		if tnow >= t0 && tnow <= t1 {
			return true
		}
	} else if t0 > t1 { // 跨天
		if (tnow >= t0 && tnow <= 24*60) || tnow <= t1 {
			return true
		}
	}

	return false
}

func Enforce(uid, tenantID uint64, obj, act string) bool {
	roles := accessctl.GetRoleForUserInDomain(uid, tenantID)
	for i := 0; i < len(roles); i++ {
		if strings.Compare(roles[i], "root") == 0 {
			return true
		}
	}

	access, _ := accessctl.Enforce(uid, tenantID, obj, act)

	return access
}
