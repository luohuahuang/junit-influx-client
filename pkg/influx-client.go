package pkg

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func GetDefaultInfluxClient() influxdb2.Client {
	return influxdb2.NewClient(InfluxUrl, InfluxToken)
}

func WriteToInflux(measurement string, tags map[string]string, fields map[string]interface{}) error {
	ctx := context.Background()
	client := GetDefaultInfluxClient()
	defer client.Close()
	writeClient := client.WriteAPIBlocking(InfluxOrg, InfluxBucket)
	p := influxdb2.NewPoint(measurement,
		tags,
		fields,
		time.Now())
	err := writeClient.WritePoint(ctx, p)
	if err != nil {
		log.Println("Error when writing to Influx: ", err, "for measurement: ", measurement)
		return err
	}
	return nil
}
