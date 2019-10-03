package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

// flags
var httpPort = flag.String("http", ":7777", "set node http port")
var neighborEndpoint = flag.String("neighbor", "http://127.0.0.1:8888", "set neighbor endpoint")
var producer = flag.Bool("producer", false, "start node as a producer")
var consumer = flag.Bool("consumer", false, "start node as a consumer")

func init() {
	flag.Parse()

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
	log.Info().Msg("starting app")

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

	log.Info().Msg("dependencies initialized")

	if *producer && *consumer {
		log.Fatal().Msg("you cannot start node as a producer and as a consumer")
	}

	switch true {
	case *producer:
	case *consumer:
	default:
		log.Fatal().Msg("you need to specify node type (consumer or producer)")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-interrupt
	log.Info().Msg("handle SIGINT, SIGTERM, SIGQUIT")
}
