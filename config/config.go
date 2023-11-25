package config

import (
	"log"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

var ValueOf = &config{}

type config struct {
	ApiID        int32  `envconfig:"API_ID" required:"true"`
	ApiHash      string `envconfig:"API_HASH" required:"true"`
	BotToken     string `envconfig:"BOT_TOKEN" required:"true"`
	LogChannelID int64  `envconfig:"LOG_CHANNEL" required:"false"`
	Dev          bool   `envconfig:"DEV" default:"false"`
	Port         int    `envconfig:"PORT" default:"8080"`
}

func (c *config) setupEnvVars() {
	err := envconfig.Process("", c)
	if err != nil {
		panic(err)
	}
}

func Load() {
	ValueOf.setupEnvVars()
	ValueOf.LogChannelID = int64(stripInt(int(ValueOf.LogChannelID)))
	log.Println("Loaded config")
}
func stripInt(a int) int {
	strA := strconv.Itoa(abs(a))
	lastDigits := strings.Replace(strA, "100", "", 1)
	result, err := strconv.Atoi(lastDigits)
	if err != nil {
		log.Println(err)
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
