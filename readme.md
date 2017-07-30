# Shorty

A simple url shortener

## Backend

Backend contains the API and a web server that can interact with the API.
Tags are stored in a redis store for fast access

### Building

Simply go to backend directory and do `make build` then start the docker container with `docker-compose up`.
You can also do `make dev`.
You can change the port in the docker-compose.yml file (default is 8080)

### Routes

- '/' and '/web/' are for the web Client ('/' redirects to '/web/')
- '/s/{tag}' redirects to the url of the corresponding tag (if tag doesn't exists it doesn't redirect)

- '/shorten' is one of the API endpoint you can send it json data like `{url: "http://myfancy.url.com"}` and it will return `{`
    `tag: "aWeSoMeTaG",`
    `url: "http://myfancy.url.com",`
    `(error: "maybe an error occured")`
    `}`

- '/unshorten/ in a similar manner provides an endpoint to unshorten urls via tags, you can send it json data like `{tag: "aWeSoMeTaG"}` and it will return the same json as /shorten


## CLI

shortcli provides a simple way of interacting with the API:
```
  -r boolean
        To get url from a tag (default is false)
  --t string
        The tag from which to rebuild the url
  --u string
        Url you want to shorten
```


## TODO

- Finish Web client (quite ugly ATM)
- Check if url that are submitted are valid urls (by trying to request for example)
- InfluxDB store with analytics
- Display analytics on the Web client and on the CLI