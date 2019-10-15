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
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const redisPingCmd = "PING"

// envs
var redisURL string
var arangodbURL string
var arangodbUser string
var arangodbPass string

// flags
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
		case "REDIS_URL":
			redisURL = envVal
		case "ARANGODB_URL":
			arangodbURL = envVal
		case "ARANGODB_USER":
			arangodbUser = envVal
		case "ARANGODB_PASS":
			arangodbPass = envVal
		}
	}
}

func main() {
	log.Info().Msg("starting app")

	// setup redis
	redisConn, err := redis.DialURL(redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to redis error")
	}
	defer redisConn.Close()

	if _, err := redisConn.Do(redisPingCmd); err != nil {
		log.Fatal().Err(err).Msg("redis ping message error")
	}

	// setup arangodb
	arangodbConn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{arangodbURL},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("connect to arangodb error")
	}
	_, err = driver.NewClient(driver.ClientConfig{
		Connection:     arangodbConn,
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
