package config

import (
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

var ValueOf = &config{}

type config struct {
	ApiID        int32  `envconfig:"API_ID" required:"true"`
	ApiHash      string `envconfig:"API_HASH" required:"true"`
	BotToken     string `envconfig:"BOT_TOKEN" required:"true"`
	LogChannelID int64  `envconfig:"LOG_CHANNEL" required:"true"`
	Dev          bool   `envconfig:"DEV" default:"false"`
	Port         int    `envconfig:"PORT" default:"8080"`
	Host         string `envconfig:"HOST" default:"http://localhost:8080"`
	HashLength   int    `envconfig:"HASH_LENGTH" default:"6"`
	UseSessionFile bool   `envconfig:"USE_SESSION_FILE" default:"true"`
	MultiTokens    []string
}

var botTokenRegex = regexp.MustCompile(`MULTI\_TOKEN[\d+]=(.*)`)

func (c *config) setupEnvVars() {
	err := envconfig.Process("", c)
	if err != nil {
		panic(err)
	}
	val := reflect.ValueOf(c).Elem()
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "MULTI_TOKEN") {
			c.MultiTokens = append(c.MultiTokens, botTokenRegex.FindStringSubmatch(env)[1])
		}
	}
	val.FieldByName("MultiTokens").Set(reflect.ValueOf(c.MultiTokens))
}

func Load(log *zap.Logger) {
	log = log.Named("Config")
	defer log.Info("Loaded config")
	ValueOf.setupEnvVars()
	ValueOf.LogChannelID = int64(stripInt(log, int(ValueOf.LogChannelID)))
	if ValueOf.HashLength == 0 {
		log.Sugar().Info("HASH_LENGTH can't be 0, defaulting to 6")
		ValueOf.HashLength = 6
	}
	if ValueOf.HashLength > 32 {
		log.Sugar().Info("HASH_LENGTH can't be more than 32, changing to 32")
		ValueOf.HashLength = 32
	}
	if ValueOf.HashLength < 5 {
		log.Sugar().Info("HASH_LENGTH can't be less than 5, defaulting to 6")
		ValueOf.HashLength = 6
	}
}

func stripInt(log *zap.Logger, a int) int {
	strA := strconv.Itoa(abs(a))
	lastDigits := strings.Replace(strA, "100", "", 1)
	result, err := strconv.Atoi(lastDigits)
	if err != nil {
		log.Sugar().Fatalln(err)
		return 0
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
