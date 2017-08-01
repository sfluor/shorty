package main

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	client "github.com/influxdata/influxdb/client/v2"
)

// Name of our database
const dbName = "shorty"

func addData(influxClient client.Client, url string) {
	// Preparing our batch points to add data
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbName,
		Precision: "s",
	})

	if err != nil {
		logrus.Errorf("Couldn't create a batch point: %s", err)
	}

	tags := map[string]string{
		"region": "France",
	}
	fields := map[string]interface{}{
		"url": url,
	}

	pt, err := client.NewPoint(
		"redirect",
		tags,
		fields,
		time.Now(),
	)

	if err != nil {
		logrus.Errorf("Couldn't create an influxdb point: %s", err)
	}

	bp.AddPoint(pt)

	// Write our batch
	if err := influxClient.Write(bp); err != nil {
		logrus.Errorf("Coudln't write batch points: %s", err)
	}
}

func queryDB(influxClient client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: dbName,
	}
	if response, err := influxClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func createDB(influxClient client.Client) {
	queryDB(influxClient, fmt.Sprintf("CREATE DATABASE %s", dbName))
}

func getDataOfUrl(influxClient client.Client, url string) interface{} {
	res, err := queryDB(
		influxClient,
		fmt.Sprintf("SELECT * FROM %s WHERE url = '%s'", "redirect", url),
	)

	if err != nil {
		logrus.Errorf("Couldn't query influxdb database: %s", err)
	}
	result := res[0].Series[0].Values
	logrus.Infof("Found data: %v", result)
	return result
}
