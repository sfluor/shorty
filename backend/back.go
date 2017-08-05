package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	client "github.com/influxdata/influxdb/client/v2"
)

// Url type to decode json from post request
type Url struct {
	Url string
}

// Shortened type to send json
type Shortened struct {
	Tag   string `json:"tag"`
	Url   string `json:"url"`
	Error string `json:"error",omitempty`
}

// Analytics type to send json
type Analytics struct {
	ClickNumber int     `json:"clickNumber"`
	ClickTimes  []int64 `json:"clickTimes"`
	Error       string  `json:"error",omitempty`
}

// Redis client
var redisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "",
	DB:       0,
})

// InfluxDB client
var influxClient, err = client.NewHTTPClient(client.HTTPConfig{
	Addr: "http://influxdb:8086",
})

func main() {
	// If connection to influxdb failed
	if err != nil {
		logrus.Fatalf("Couldn't connect to Influxdb: %s", err)
	}
	logrus.Info("Connected to Influxdb")
	createDB(influxClient)

	pong, err := redisClient.Ping().Result()
	if err != nil {
		logrus.Fatalf("Couldn't connect to redis store: %s", err)
	}
	logrus.Infof("Connected to the redis store: %s", pong)

	r := mux.NewRouter()

	// Web app
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("static"))))

	// Redirect everything from '/' to /web/
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", 302)
	})

	// Path to shorten links
	r.HandleFunc("/shorten", shortener).Methods("POST")
	r.HandleFunc("/unshorten", unShortener)

	// Standard path to redirect
	r.HandleFunc("/s/{token}", redirect)

	// Analytics
	r.HandleFunc("/analytics/{token}", analytics)

	http.Handle("/", r)
	logrus.Info("Server started")

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func shortener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	handleInternalError(w, err)

	var u Url
	err = json.Unmarshal(b, &u)
	handleInternalError(w, err)

	logrus.Infof("Asking to shorten: %s", u.Url)
	tag := shorten(redisClient, u.Url)
	logrus.Infof("Result is tag: %s", tag)
	// Return JSON
	payload := Shortened{Tag: tag, Url: u.Url}
	err = json.NewEncoder(w).Encode(payload)
	handleInternalError(w, err)
}

func unShortener(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	logrus.Infof("This is a tag: %s", tag)

	// Check if tag is not empty
	if tag != "" {
		// TODO: check if url is tag exists if it doesn't return an error
		url := unShorten(redisClient, tag)
		payload := Shortened{Tag: tag, Url: url}
		// Return JSON
		err := json.NewEncoder(w).Encode(payload)
		handleInternalError(w, err)
	} else {
		fmt.Fprintf(w, "Sorry your tag is empty")
	}

}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// We parse the url otherwise url likes google.com redirect internally
	url, err := formatUrl(unShorten(redisClient, vars["token"]))
	handleInternalError(w, err)

	addData(influxClient, url)
	http.Redirect(w, r, url, 302)
	fmt.Fprint(w, "Redirecting...")
}

// formatUrl correctly formats url
//	Example:
// 		formatUrl("google.com") outputs "http://google.com"
func formatUrl(data string) (string, error) {
	_url, err := url.Parse(data)
	// If scheme is empty (ie not http or https) just put http
	if _url.Scheme == "" {
		return fmt.Sprintf("http://%s", _url.String()), nil
	}
	// Else return the url
	return _url.String(), err
}

func analytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	url, err := formatUrl(unShorten(redisClient, vars["token"]))
	handleInternalError(w, err)
	data, ok := getDataOfUrl(influxClient, url).([][]interface{})

	// Type assertion failed
	if !ok {
		handleInternalError(w, errors.New("Type assertion failed"))
	}

	// Extract our data
	fData, err := extractTime(data)
	handleInternalError(w, err)

	payload := Analytics{ClickNumber: len(fData), ClickTimes: fData}
	err = json.NewEncoder(w).Encode(payload)
	handleInternalError(w, err)
}

func handleInternalError(w http.ResponseWriter, err error) {
	if err != nil {
		logrus.Errorf("An error occured: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
