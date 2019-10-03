package main

import (
	"os"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// envs
var natsUrl string
var natsUser string
var natsPass string
var arangodbUrl string
var arangodbUser string
var arangodbPass string

func init() {
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	for _, envPair := range os.Environ() {
		splittedEnvPair := strings.Split(envPair, "=")
		envKey, envVal := splittedEnvPair[0], splittedEnvPair[1]
		switch envKey {
		case "NATS_URL":
			natsUrl = envVal
		case "NATS_USER":
			natsUser = envVal
		case "NATS_PASS":
			natsPass = envVal
		case "ARANGODB_URL":
			arangodbUrl = envVal
		case "ARANGODB_USER":
			arangodbUser = envVal
		case "ARANGODB_PASS":
			arangodbPass = envVal
		}
	}
}

func main() {
	// setup nats
	natsConn, err := nats.Connect(natsUrl, nats.UserInfo(natsUser, natsPass))
	if err != nil {
		log.Fatal().Err(err).Msg("connect to nats error")
	}
	defer natsConn.Close()

	// setup arangodb
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{arangodbUrl},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("connect to arangodb error")
	}
	_, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(arangodbUser, arangodbPass),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("initialize arangodb client error")
	}
}
