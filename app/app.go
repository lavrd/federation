package main

import (
	"context"
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
const redisPublishCmd = "PUBLISH"

const redisTopicName = "default"

const arangodbDatabaseName = "default"
const arangodbCollectionName = "default"

// envs
var redisURL string
var arangodbURL string
var arangodbUser string
var arangodbPass string

// flags
var startAsProducer = flag.Bool("producer", false, "start node as a producer")
var startAsConsumer = flag.Bool("consumer", false, "start node as a consumer")

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
	rconn, err := redis.DialURL(redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to redis error")
	}
	defer rconn.Close()

	if _, err := rconn.Do(redisPingCmd); err != nil {
		log.Fatal().Err(err).Msg("redis ping message error")
	}

	// setup arangodb
	aconn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{arangodbURL},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("connect to arangodb error")
	}
	aclient, err := driver.NewClient(driver.ClientConfig{
		Connection:     aconn,
		Authentication: driver.BasicAuthentication(arangodbUser, arangodbPass),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("initialize arangodb client error")
	}

	log.Info().Msg("dependencies initialized")

	if *startAsProducer && *startAsConsumer {
		log.Fatal().Msg("you cannot start node as a producer and as a consumer")
	}

	switch true {
	case *startAsProducer:
		go func() {
			if err := runProducer(rconn); err != nil {
				log.Fatal().Err(err).Msg("run producer error")
			}
		}()
	case *startAsConsumer:
		if err := runConsumer(rconn, aclient); err != nil {
			log.Fatal().Err(err).Msg("run consumer error")
		}
	default:
		log.Fatal().Msg("you need to specify node type (consumer or producer)")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-interrupt
	log.Info().Msg("handle SIGINT, SIGTERM, SIGQUIT")
}

func runProducer(rconn redis.Conn) error {
	for range time.NewTicker(time.Second).C {
		_, err := rconn.Do(redisPublishCmd, redisTopicName, "value")
		if err != nil {
			return err
		}
	}
	return nil
}

func runConsumer(rconn redis.Conn, aclient driver.Client) error {
	rpsconn := redis.PubSubConn{Conn: rconn}

	err := rpsconn.Subscribe(redisTopicName)
	if err != nil {
		return err
	}

	exists, err := aclient.DatabaseExists(context.Background(), arangodbDatabaseName)
	if err != nil {
		return err
	}
	if !exists {
		_, err := aclient.CreateDatabase(context.Background(), arangodbDatabaseName, nil)
		if err != nil {
			return err
		}
	}
	db, err := aclient.Database(context.Background(), arangodbDatabaseName)
	if err != nil {
		return err
	}

	exists, err = db.CollectionExists(context.Background(), arangodbCollectionName)
	if err != nil {
		return err
	}
	if !exists {
		_, err := db.CreateCollection(context.Background(), arangodbCollectionName, nil)
		if err != nil {
			return err
		}
	}
	col, err := db.Collection(context.Background(), arangodbCollectionName)
	if err != nil {
		return err
	}

	for {
		switch recv := rpsconn.Receive().(type) {
		case redis.Message:
			log.Info().Msgf("`%s` ch received `%s`", recv.Channel, recv.Data)
			_, err := col.CreateDocument(context.Background(), struct {
				Data string `json:"data"`
			}{
				Data: string(recv.Data),
			})
			if err != nil {
				return err
			}
		case error:
			return recv
		}
	}
}
