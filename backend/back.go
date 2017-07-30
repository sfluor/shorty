package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Url type to decode json from post request
type Url struct {
	Url string
}

// Shortened type to send json
type Shortened struct {
	Tag   string `json:"tag"`
	Url   string `json:"url"`
	Error string `json:error,omitempty`
}

var client *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "",
	DB:       0,
})

func main() {
	pong, err := client.Ping().Result()
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
	http.Handle("/", r)
	logrus.Info("Server started")

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func shortener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("An error occured: %s", err)
	}
	var u Url
	if err := json.Unmarshal(b, &u); err != nil {
		logrus.Errorf("An error occured during json decoding: %s", err)
	}
	logrus.Infof("Asking to shorten: %s", u.Url)
	tag := shorten(client, u.Url)
	logrus.Infof("Result is tag: %s", tag)
	// Return JSON
	payload := Shortened{Tag: tag, Url: u.Url}
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		logrus.Errorf("An error occured during json encoding: %s", err)
	}
}

func unShortener(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	logrus.Infof("This is a tag: %s", tag)

	// Check if tag is not empty
	if tag != "" {
		// TODO: check if url is tag exists if it doesn't return an error
		url := unShorten(client, tag)
		payload := Shortened{Tag: tag, Url: url}
		// Return JSON
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			logrus.Errorf("An error occured during json encoding: %s", err)
		}
	} else {
		fmt.Fprintf(w, "Sorry your tag is empty")
	}

}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// We parse the url otherwise url likes google.com redirect internally
	url := formatUrl(unShorten(client, vars["token"]))

	http.Redirect(w, r, url, 302)
	fmt.Fprint(w, "Redirecting...")
}

// formatUrl correctly formats url
//	Example:
// 		formatUrl("google.com") outputs "http://google.com"
func formatUrl(data string) string {
	_url, err := url.Parse(data)
	if err != nil {
		logrus.Errorf("An error occured when parsing the url: %s", err)
	}
	// If scheme is empty (ie not http or https) just put http
	if _url.Scheme == "" {
		return fmt.Sprintf("http://%s", _url.String())
	}
	// Else return the url
	return _url.String()
}
