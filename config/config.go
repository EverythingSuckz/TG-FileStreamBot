package config

import (
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
}

func (c *config) setupEnvVars() {
	err := envconfig.Process("", c)
	if err != nil {
		panic(err)
	}
}

func Load(log *zap.Logger) {
	ValueOf.setupEnvVars()
	ValueOf.LogChannelID = int64(stripInt(log, int(ValueOf.LogChannelID)))
	log.Info("Loaded config")
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
