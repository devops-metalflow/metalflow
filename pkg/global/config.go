package global

import "go.uber.org/zap/zapcore"

// Configuration 系统配置, 配置字段可参见yml注释
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
type Configuration struct {
	System    SystemConfiguration    `mapstructure:"system" json:"system"`
	Logs      LogsConfiguration      `mapstructure:"logs" json:"logs"`
	Mysql     MysqlConfiguration     `mapstructure:"mysql" json:"mysql"`
	Redis     RedisConfiguration     `mapstructure:"redis" json:"redis"`
	Jwt       JwtConfiguration       `mapstructure:"jwt" json:"jwt"`
	RateLimit RateLimitConfiguration `mapstructure:"rate-limit" json:"rateLimit"`
	Consul    ConsulConfiguration    `mapstructure:"consul" json:"consul"`
	Upload    UploadConfiguration    `mapstructure:"upload" json:"upload"`
	NodeConf  NodeConfiguration      `mapstructure:"node" json:"node"`
	Mail      MailConfiguration      `mapstructure:"mail" json:"mail"`
}

type SystemConfiguration struct {
	UrlPathPrefix               string   `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	ApiVersion                  string   `mapstructure:"api-version" json:"apiVersion"`
	Port                        int      `mapstructure:"port" json:"port"`
	PprofPort                   int      `mapstructure:"pprof-port" json:"pprofPort"`
	ConnectTimeout              int      `mapstructure:"connect-timeout" json:"connectTimeout"`
	ExecuteTimeout              int      `mapstructure:"execute-timeout" json:"executeTimeout"`
	Transaction                 bool     `mapstructure:"transaction" json:"transaction"`
	InitData                    bool     `mapstructure:"init-data" json:"initData"`
	OperationLogKey             string   `mapstructure:"operation-log-key" json:"operationLogKey"`
	OperationLogDisabledPaths   string   `mapstructure:"operation-log-disabled-paths" json:"operationLogDisabledPaths"`
	OperationLogDisabledPathArr []string `mapstructure:"-" json:"-"`
	OperationLogAllowedToDelete bool     `mapstructure:"operation-log-allowed-to-delete" json:"operationLogAllowedToDelete"`
	IdempotenceTokenName        string   `mapstructure:"idempotence-token-name" json:"idempotenceTokenName"`
	NodeMetricsCronTask         string   `mapstructure:"node-metrics-cron-task" json:"nodeMetricsCronTask"`
	NodePingCronTask            string   `mapstructure:"node-ping-cron-task" json:"nodePingCronTask"`
}

type LogsConfiguration struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type NodeConfiguration struct {
	AddrBind []NodeAddrConfiguration `mapstructure:"addr-bind" json:"addrBind"`
	Hide     string                  `mapstructure:"hide" json:"hide"`
}

type NodeAddrConfiguration struct {
	Addr string   `mapstructure:"addr" json:"addr"`
	Ips  []string `mapstructure:"ips" json:"ips"`
}

type MysqlConfiguration struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}

type RedisConfiguration struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      int    `mapstructure:"port" json:"port"`
	Password  string `mapstructure:"password" json:"password"`
	Database  int    `mapstructure:"database" json:"database"`
	BinlogPos string `mapstructure:"binlog-pos" json:"binlogPos"`
}

type ConsulConfiguration struct {
	Address string `mapstructure:"address" json:"address"`
	Port    int    `mapstructure:"port" json:"port"`
}

type JwtConfiguration struct {
	Realm      string `mapstructure:"realm" json:"realm"`
	Key        string `mapstructure:"key" json:"key"`
	Timeout    int    `mapstructure:"timeout" json:"timeout"`
	MaxRefresh int    `mapstructure:"max-refresh" json:"maxRefresh"`
}

type RateLimitConfiguration struct {
	Max int64 `mapstructure:"max" json:"max"`
}

type UploadConfiguration struct {
	SaveDir              string `mapstructure:"save-dir" json:"saveDir"`
	SingleMaxSize        uint   `mapstructure:"single-max-size" json:"singleMaxSize"`
	MergeConcurrentCount uint   `mapstructure:"merge-concurrent-count" json:"mergeConcurrentCount"`
}

type MailConfiguration struct {
	Host     string   `mapstructure:"host" json:"host"`
	Port     int      `mapstructure:"port" json:"port"`
	Username string   `mapstructure:"username" json:"username"`
	Password string   `mapstructure:"password" json:"password"`
	From     string   `mapstructure:"from" json:"from"`
	Header   string   `mapstructure:"header" json:"header"`
	Suffix   string   `mapstructure:"suffix" json:"suffix"`
	Cc       []string `mapstructure:"cc" json:"cc"`
}
