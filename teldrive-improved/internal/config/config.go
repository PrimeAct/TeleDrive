package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toKebabCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func getKey(f reflect.StructField) string {
	if t := f.Tag.Get("koanf"); t != "" {
		return t
	}
	return toKebabCase(f.Name)
}

type ServerHTTPConfig struct {
	Host         string        `default:"0.0.0.0" description:"Server host address"`
	Port         int           `default:"8080" description:"Server port"`
	AutoPort     bool          `default:"false" description:"Automatically find available port"`
	ReadTimeout  time.Duration `default:"30s" description:"HTTP read timeout"`
	WriteTimeout time.Duration `default:"30s" description:"HTTP write timeout"`
	IdleTimeout  time.Duration `default:"120s" description:"HTTP idle timeout"`
	CORS         CORSConfig    `koanf:"cors"`
}

type CORSConfig struct {
	Enabled          bool          `default:"true" description:"Enable CORS"`
	AllowedOrigins   []string      `default:"*" description:"Allowed CORS origins"`
	AllowedMethods   []string      `default:"GET,POST,PUT,DELETE,OPTIONS,PATCH" description:"Allowed HTTP methods"`
	AllowedHeaders   []string      `default:"*" description:"Allowed HTTP headers"`
	ExposedHeaders   []string      `default:"Content-Length,Content-Range" description:"Exposed HTTP headers"`
	AllowCredentials bool          `default:"true" description:"Allow credentials in CORS"`
	MaxAge           time.Duration `default:"24h" description:"CORS max age"`
}

type LogConfig struct {
	Level      string `default:"info" description:"Log level (debug, info, warn, error, fatal, panic)"`
	Format     string `default:"json" description:"Log format (json, console)"`
	Output     string `default:"stdout" description:"Log output path (stdout, stderr, or file path)"`
	MaxSize    int    `default:"100" description:"Max log file size in MB before rotation"`
	MaxBackups int    `default:"5" description:"Max number of old log files to retain"`
	MaxAge     int    `default:"30" description:"Max days to retain old log files"`
	Compress   bool   `default:"true" description:"Compress rotated log files"`
}

type JWTConfig struct {
	Secret       string        `default:"" description:"JWT signing secret (auto-generated if empty)"`
	SessionTime  time.Duration `default:"720h" description:"JWT session duration"`
	AllowedUsers []string      `default:"" description:"List of allowed Telegram usernames (empty = all)"`
}

type DBConfig struct {
	Host    string     `default:"localhost" description:"Database host"`
	Port    int        `default:"5432" description:"Database port"`
	User    string     `default:"postgres" description:"Database user"`
	Password string    `default:"postgres" description:"Database password"`
	Name    string     `default:"teldrive" description:"Database name"`
	SSLMode string     `default:"disable" description:"SSL mode (disable, require, verify-ca, verify-full)"`
	Pool    PoolConfig `koanf:"pool"`
}

type PoolConfig struct {
	MaxOpen     int           `default:"25" description:"Max open connections"`
	MaxIdle     int           `default:"5" description:"Max idle connections"`
	MaxLifetime time.Duration `default:"1h" description:"Max connection lifetime"`
	MaxIdleTime time.Duration `default:"30m" description:"Max idle connection time"`
}

type TGConfig struct {
	AppID            int64       `default:"" description:"Telegram App ID"`
	AppHash          string      `default:"" description:"Telegram App Hash"`
	Phone            string      `default:"" description:"Phone number (with country code)"`
	Password         string      `default:"" description:"2FA password if enabled"`
	SessionDir       string      `default:"./sessions" description:"Directory for session files"`
	Workers          int         `default:"4" description:"Number of upload/download workers"`
	MultiBots        []string    `default:"" description:"List of bot tokens for multi-bot upload"`
	UploadPoolSize   int         `default:"50" description:"Upload pool size"`
	DownloadPoolSize int         `default:"50" description:"Download pool size"`
	Proxy            ProxyConfig `koanf:"proxy"`
}

type ProxyConfig struct {
	Enabled  bool   `default:"false" description:"Enable proxy"`
	URL      string `default:"" description:"Proxy URL (socks5:// or http://)"`
	Username string `default:"" description:"Proxy username"`
	Password string `default:"" description:"Proxy password"`
}

type CronJobConfig struct {
	Enabled         bool          `default:"true" description:"Enable cron jobs"`
	CleanupInterval time.Duration `default:"24h" description:"Cleanup job interval"`
	StatsInterval   time.Duration `default:"1h" description:"Stats collection interval"`
}

type CacheConfig struct {
	Type    string        `default:"memory" description:"Cache type (memory, redis)"`
	TTL     time.Duration `default:"5m" description:"Default cache TTL"`
	MaxSize int           `default:"10000" description:"Max cache entries (memory only)"`
}

type RedisConfig struct {
	Enabled  bool   `default:"false" description:"Enable Redis"`
	Address  string `default:"localhost:6379" description:"Redis server address"`
	Password string `default:"" description:"Redis password"`
	DB       int    `default:"0" description:"Redis database number"`
}

type EventConfig struct {
	PollInterval     time.Duration `default:"10s" description:"Event polling interval"`
	DBWorkers        int           `default:"10" description:"Number of DB worker goroutines"`
	DBBufferSize     int           `default:"1000" description:"DB worker queue buffer size"`
	DeduplicationTTL time.Duration `default:"5s" description:"Event deduplication TTL"`
}

type ServerConfig struct {
	Server   ServerHTTPConfig `koanf:"server"`
	Log      LogConfig        `koanf:"log"`
	JWT      JWTConfig        `koanf:"jwt"`
	DB       DBConfig         `koanf:"db"`
	TG       TGConfig         `koanf:"tg"`
	CronJobs CronJobConfig    `koanf:"cron-jobs"`
	Cache    CacheConfig      `koanf:"cache"`
	Redis    RedisConfig      `koanf:"redis"`
	Events   EventConfig      `koanf:"events"`
}

type CheckConfig struct {
	Log        LogConfig `koanf:"log"`
	DB         DBConfig  `koanf:"db"`
	TG         TGConfig  `koanf:"tg"`
	ExportFile string    `default:"results.json" description:"Path for exported JSON file"`
	DryRun     bool      `default:"false" description:"Simulate without making changes"`
	User       string    `default:"" description:"Telegram username to check"`
	Concurrent int       `default:"4" description:"Number of concurrent workers"`
}

type Loader struct {
	k *koanf.Koanf
}

func NewLoader() *Loader {
	return &Loader{k: koanf.New(".")}
}

func (l *Loader) Load(cmd *cobra.Command, cfg interface{}) error {
	configPaths := []string{
		"config.toml",
		"config.yml",
		"config.yaml",
		filepath.Join(os.Getenv("HOME"), ".config", "teldrive", "config.toml"),
		"/etc/teldrive/config.toml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		var parser koanf.Parser
		switch filepath.Ext(path) {
		case ".toml":
			parser = toml.Parser()
		case ".yml", ".yaml":
			parser = yaml.Parser()
		default:
			continue
		}

		if err := l.k.Load(file.Provider(path), parser); err != nil {
			return fmt.Errorf("failed to load config from %s: %w", path, err)
		}
	}

	if err := l.k.Load(env.Provider("TELDRIVE_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "TELDRIVE_")), "_", ".")
	}), nil); err != nil {
		return fmt.Errorf("failed to load env config: %w", err)
	}

	if cmd != nil {
		if err := l.k.Load(posflagProvider(cmd.Flags(), ".", l.k), nil); err != nil {
			return fmt.Errorf("failed to load flag config: %w", err)
		}
	}

	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		Result:           cfg,
		WeaklyTypedInput: true,
		TagName:          "koanf",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(l.k.All()); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	return nil
}

func (l *Loader) Validate(cfg interface{}) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func (l *Loader) RegisterFlags(fs *pflag.FlagSet, t reflect.Type, prefix ...string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("skipPflag") == "true" {
			continue
		}

		key := getKey(field)
		if len(prefix) > 0 {
			key = prefix[0] + "." + key
		}

		if field.Type.Kind() == reflect.Struct {
			l.RegisterFlags(fs, field.Type, key)
			continue
		}

		defaultVal := field.Tag.Get("default")
		desc := field.Tag.Get("description")
		flagName := strings.ReplaceAll(key, ".", "-")

		switch field.Type.Kind() {
		case reflect.String:
			fs.String(flagName, defaultVal, desc)
		case reflect.Int, reflect.Int64:
			if v, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
				fs.Int64(flagName, v, desc)
			}
		case reflect.Bool:
			if v, err := strconv.ParseBool(defaultVal); err == nil {
				fs.Bool(flagName, v, desc)
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				fs.StringSlice(flagName, nil, desc)
			}
		}
	}
}

func posflagProvider(fs *pflag.FlagSet, delim string, k *koanf.Koanf) *koanf.Pflag {
	// Simplified - in real implementation use proper koanf pflag provider
	return nil
}
