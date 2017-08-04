package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// Shortened type to send json
type Shortened struct {
	Tag   string `json:"tag"`
	Url   string `json:"url"`
	Error string `json:error,omitempty`
}

// Analytics type to receive json
type Analytics struct {
	ClickNumber int `json:"clickNumber"`
}

const backendUrl = "http://localhost:8080"

func main() {
	url := flag.String("u", "", "Url you want to shorten")
	reverse := flag.Bool("r", false, "To get url from a tag")
	tag := flag.String("t", "", "The tag from which to rebuild the url")

	flag.Parse()

	// If the user wants to reverse the url
	if *reverse && *tag != "" {
		data := unShortenURL(*tag)
		analytics := getAnalytics(*tag)
		if data.Url == "" {
			fmt.Println("Sorry this url doesn't seem to exist in records")
			return
		}
		fmt.Printf("%s/s/%s corresponds to the url: %s\n", backendUrl, *tag, data.Url)
		// Also print analytics
		fmt.Printf("Number of clicks: %v\n", analytics.ClickNumber)
	} else if *url != "" {
		data := shortenURL(*url)
		fmt.Printf("Url %s has been shortened to %s/s/%s\n", *url, backendUrl, data.Tag)
	} else {
		fmt.Print("Sorry bad input check for shorty help with command shorty -h")
	}
}

// shortenURL queries for our backend and returns the shortened url
func shortenURL(url string) Shortened {
	// Our url containing the querystring
	reqUrl := fmt.Sprintf("%s/shorten", backendUrl)

	// Data Url to send for shortening
	payload := map[string]string{"url": url}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Fatalf("An error occured during json encoding: %s", err)
	}

	// Send post Request
	resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logrus.Fatalf("An error occured during the request: %s", err)
	}
	defer resp.Body.Close()
	return extractAndSend(resp.Body)
}

// unShortenUrl queries for our backend and returns the true url
func unShortenURL(tag string) Shortened {
	// Our url containing the querystring
	reqUrl := fmt.Sprintf("%s/unshorten?tag=%s", backendUrl, tag)
	resp, err := http.Get(reqUrl)
	if err != nil {
		logrus.Fatalf("An error occured during the request: %s", err)
	}
	defer resp.Body.Close()
	return extractAndSend(resp.Body)
}

// extractAndSend get json data from the body and return it
func extractAndSend(body io.ReadCloser) Shortened {
	var target Shortened
	err := json.NewDecoder(body).Decode(&target)
	if err != nil {
		logrus.Fatalf("An error occured during the json decoding: %s", err)
	}
	return target
}

// getAnalytics simply query the analytics api
func getAnalytics(tag string) Analytics {
	reqUrl := fmt.Sprintf("%s/analytics/%s", backendUrl, tag)
	resp, err := http.Get(reqUrl)
	if err != nil {
		logrus.Fatalf("An error occured during the request: %s", err)
	}
	defer resp.Body.Close()
	var a Analytics
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		logrus.Fatalf("An error occured during the json decoding: %s", err)
	}
	return a
}
