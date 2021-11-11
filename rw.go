package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/influxdata/influxdb-client-go/v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	influxdbHost := os.Getenv("INFLUXDB_HOST")
	influxdbPort := os.Getenv("INFLUXDB_PORT")
	influxdbToken := os.Getenv("INFLUXDB_TOKEN")
	//influxdbOrg := os.Getenv("INFLUXDB_ORG")
	influxdbBucket := os.Getenv("INFLUXDB_BUCKET_NAME")
	defaultDuration := os.Getenv("INFLUX_DEFAULT_LOAD_DURATION")

	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxdbHost, influxdbPort), influxdbToken)

	// Get query client
	queryAPI := client.QueryAPI("Sea Creeper")
	// Get parser flux query result
	result, err := queryAPI.Query(context.Background(),
		`from(bucket: "`+influxdbBucket+`")
				|> range(start: -`+defaultDuration+`)
  				|> filter(fn: (r) => r["_measurement"] == "newman-result")
  				|> keep(columns: ["_time", "url"])
  				|> limit(n:10)
  				|> yield()`)
	if err == nil {
		for result.Next() {
			fmt.Printf("%s | %s\n", result.Record().Time(), result.Record().ValueByKey("url"))
		}
		if result.Err() != nil {
			fmt.Printf("Query error: %s\n", result.Err().Error())
		}
	}
	client.Close()
}
