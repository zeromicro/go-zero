package main

import (
	"flag"
	"time"

	"zero/core/conf"
	"zero/core/lang"
	"zero/core/proc"
	"zero/kq"
	"zero/stash/config"
	"zero/stash/es"
	"zero/stash/filter"
	"zero/stash/handler"

	"github.com/olivere/elastic"
)

const dateFormat = "2006.01.02"

var configFile = flag.String("f", "etc/config.json", "Specify the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	proc.SetTimeoutToForceQuit(c.GracePeriod)

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(c.Output.ElasticSearch.Hosts...),
	)
	lang.Must(err)

	indexFormat := c.Output.ElasticSearch.DailyIndexPrefix + dateFormat
	var loc *time.Location
	if len(c.Output.ElasticSearch.TimeZone) > 0 {
		loc, err = time.LoadLocation(c.Output.ElasticSearch.TimeZone)
		lang.Must(err)
	} else {
		loc = time.Local
	}
	indexer := es.NewIndex(client, func(t time.Time) string {
		return t.In(loc).Format(indexFormat)
	})

	filters := filter.CreateFilters(c)
	writer, err := es.NewWriter(c.Output.ElasticSearch, indexer)
	lang.Must(err)

	handle := handler.NewHandler(writer)
	handle.AddFilters(filters...)
	handle.AddFilters(filter.AddUriFieldFilter("url", "uri"))
	q := kq.MustNewQueue(c.Input.Kafka, handle)
	q.Start()
}
