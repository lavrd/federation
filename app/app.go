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

const documentDefaultKey = "1"

// envs
var redisURL string
var arangodbURL string
var arangodbUser string
var arangodbPass string

// flags
var startAsProducer = flag.Bool("producer", false, "start node as a producer")
var startAsConsumer = flag.Bool("consumer", false, "start node as a consumer")
var startAsValidator = flag.Bool("validator", false, "start node as a validator")

type Data struct {
	Key   string    `json:"_key"`
	Value time.Time `json:"value"`
}

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

	switch true {
	case *startAsProducer:
		go func() {
			rconn, err := initRedisConn()
			if err != nil {
				log.Fatal().Err(err).Msg("init redis connection error")
			}
			defer rconn.Close()

			if err := runProducer(rconn); err != nil {
				log.Fatal().Err(err).Msg("run producer error")
			}
		}()
	case *startAsConsumer:
		rconn, err := initRedisConn()
		if err != nil {
			log.Fatal().Err(err).Msg("init redis connection error")
		}
		defer rconn.Close()

		aclient, err := initArangodbConn()
		if err != nil {
			log.Fatal().Err(err).Msg("init arangodb connection error")
		}

		if err := runConsumer(rconn, aclient); err != nil {
			log.Fatal().Err(err).Msg("run consumer error")
		}
	case *startAsValidator:
		aclient, err := initArangodbConn()
		if err != nil {
			log.Fatal().Err(err).Msg("init arangodb connection error")
		}

		if err := runValidator(aclient); err != nil {
			log.Fatal().Err(err).Msg("run validator error")
		}
	default:
		log.Fatal().Msg("you need to specify node type (consumer or producer)")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-interrupt
	log.Info().Msg("handle SIGINT, SIGTERM, SIGQUIT")
}

func initArangodbConn() (driver.Client, error) {
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

	return aclient, nil
}

func initRedisConn() (redis.Conn, error) {
	rconn, err := redis.DialURL(redisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to redis error")
	}

	if _, err := rconn.Do(redisPingCmd); err != nil {
		log.Fatal().Err(err).Msg("redis ping message error")
	}

	return rconn, nil
}

func runProducer(rconn redis.Conn) error {
	for curTime := range time.NewTicker(time.Second).C {
		_, err := rconn.Do(redisPublishCmd, redisTopicName, curTime.Format(time.RFC3339))
		if err != nil {
			return err
		}
	}
	return nil
}

func runValidator(aclient driver.Client) error {
	col, err := initArangodbCollection(aclient)
	if err != nil {
		return err
	}

	prevTime := time.Time{}

	for range time.NewTicker(time.Second).C {
		data := &Data{}
		_, err := col.ReadDocument(context.Background(), documentDefaultKey, data)
		if driver.IsNotFound(err) {
			log.Info().Msg("document not found")
			continue
		}
		if err != nil {
			return err
		}

		if data.Value.Before(prevTime) || data.Value.Equal(prevTime) {
			log.Info().Msg("incorrect received time from database")
		} else {
			log.Info().Msg("correct received time from database")
		}

		prevTime = data.Value
	}

	return nil
}

func initArangodbCollection(aclient driver.Client) (driver.Collection, error) {
	exists, err := aclient.DatabaseExists(context.Background(), arangodbDatabaseName)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := aclient.CreateDatabase(context.Background(), arangodbDatabaseName, nil)
		if err != nil {
			return nil, err
		}
	}
	db, err := aclient.Database(context.Background(), arangodbDatabaseName)
	if err != nil {
		return nil, err
	}

	exists, err = db.CollectionExists(context.Background(), arangodbCollectionName)
	if err != nil {
		return nil, err
	}
	if !exists {
		_, err := db.CreateCollection(context.Background(), arangodbCollectionName, nil)
		if err != nil {
			return nil, err
		}
	}
	col, err := db.Collection(context.Background(), arangodbCollectionName)
	if err != nil {
		return nil, err
	}

	return col, nil
}

func runConsumer(rconn redis.Conn, aclient driver.Client) error {
	rpsconn := redis.PubSubConn{Conn: rconn}

	err := rpsconn.Subscribe(redisTopicName)
	if err != nil {
		return err
	}

	col, err := initArangodbCollection(aclient)
	if err != nil {
		return err
	}

	_, err = col.CreateDocument(context.Background(), &Data{Key: documentDefaultKey})
	if err != nil && !driver.IsConflict(err) {
		return err
	}

	for {
		switch recv := rpsconn.Receive().(type) {
		case redis.Message:
			log.Info().Msgf("`%s` ch received `%s`", recv.Channel, recv.Data)
			recvTime, err := time.Parse(time.RFC3339, string(recv.Data))
			if err != nil {
				return err
			}

			_, err = col.UpdateDocument(
				context.Background(),
				documentDefaultKey,
				&Data{Key: documentDefaultKey, Value: recvTime},
			)
			if err != nil {
				return err
			}
		case error:
			return recv
		}
	}
}
