package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var scribeName string
var influxdbHost string
var influxdbPort string
var influxdbToken string
var influxdbOrg string
var influxdbBucket string
var influxdbMeasurement string
var influxdbDefaultDuration string
var influxdbDefaultLimit string

type scribe struct {
	Time string `json:"time"`
	Url  string `json:"url"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Errorln("No .env file found")
	}
	scribeName = os.Getenv("SCRIBER_NAME")
	influxdbHost = os.Getenv("INFLUXDB_HOST")
	influxdbPort = os.Getenv("INFLUXDB_PORT")
	influxdbToken = os.Getenv("INFLUXDB_TOKEN")
	influxdbOrg = os.Getenv("INFLUXDB_ORG")
	influxdbBucket = os.Getenv("INFLUXDB_BUCKET_NAME")
	influxdbMeasurement = os.Getenv("INFLUXDB_MEASUREMENT")
	influxdbDefaultDuration = os.Getenv("INFLUX_DEFAULT_LOAD_DURATION")
	influxdbDefaultLimit = os.Getenv("INFLUX_DEFAULT_LIMIT")

	log.Debugln("Scribe name: " + scribeName)
	log.Debugln("InfluxDB measurement: " + influxdbMeasurement)
	server := gin.Default()
	server.GET("/:scriber", func(c *gin.Context) {
		scribe := c.Param("scriber")
		if scribe != scribeName {
			serveErr(c)
		} else {
			serve(c)
		}
	})
	server.Run()
}

func serveErr(c *gin.Context) {
	status := http.StatusMethodNotAllowed
	c.Header("Allow", scribeName)
	c.JSON(status, gin.H{
		"code":    status,
		"message": http.StatusText(status),
	})
}

func query(start string, stop string, limit string) *api.QueryTableResult {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxdbHost, influxdbPort), influxdbToken)

	// Get query client
	queryAPI := client.QueryAPI(influxdbOrg)
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(),
		`from(bucket: "`+influxdbBucket+`")
				|> range(start: `+start+`, stop: `+stop+`)
  				|> filter(fn: (r) => r["_measurement"] == "`+influxdbMeasurement+`")
  				|> keep(columns: ["_time", "url"])
				|> sort(desc: true)
  				|> limit(n:`+limit+`)
  				|> yield()`)
	if err == nil {
		if result.Err() != nil {
			log.Errorln("Query error: %s", result.Err().Error())
		}
	} else {
		log.Errorln("InfluxDB client error: ", err.Error())
	}
	client.Close()
	return result
}

func serve(c *gin.Context) {
	/**
	TODO Prevent query injection by doing client side validation
		 Parametrized Flux queries are not yet supported
	     See https://github.com/influxdata/influxdb-client-go/issues/146
	*/
	start := c.DefaultQuery("start", influxdbDefaultDuration)
	stop := c.DefaultQuery("stop", "now()")
	limit := c.DefaultQuery("limit", influxdbDefaultLimit)
	response := query(start, stop, limit)
	var result []scribe
	for response.Next() {
		result = append(result, scribe{Time: response.Record().Time().String(), Url: response.Record().ValueByKey("url").(string)})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
