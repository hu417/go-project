package logger

import (
	"os"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
)

func Test_zerolog(t *testing.T) {

	// ==============================================================================================
	//
	// Set up a logger
	//
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	// ==============================================================================================
	//
	// Pass the logger to the client
	//
	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Logger: &CustomLogger{log},
	})

	// ----------------------------------------------------------------------------------------------
	{
		es.Delete("test", "1")
		es.Exists("test", "1")
		es.Index("test", strings.NewReader(`{"title" : "logging"}`), es.Index.WithRefresh("true"))

		es.Search(
			es.Search.WithQuery("{FAIL"),
		)

		es.Search(
			es.Search.WithIndex("test"),
			es.Search.WithBody(strings.NewReader(`{"query" : {"match" : { "title" : "logging" } } }`)),
			es.Search.WithSize(1),
		)
	}
}
