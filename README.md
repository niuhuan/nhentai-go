NHENTAI-GO
================

nHentai api with golang

## Features

- List or search comics
- Get comic info and comic page images

## Usage

```go
package main

// import the package
import "github.com/niuhuan/nhentai-go"
import "net/http"
import "net/url"
import "time"

func main() {
	// new client
	client := nhentai.Client{}
	// set proxy (optional)
	proxyUrl, proxyErr := url.Parse("socks5://127.0.0.1:1080")
	proxy := func(_ *http.Request) (*url.URL, error) {
		return proxyUrl, proxyErr
	}
	client.Transport = &http.Transport{
		Proxy:                 proxy,
		TLSHandshakeTimeout:   time.Second * 10,
		ExpectContinueTimeout: time.Second * 10,
		ResponseHeaderTimeout: time.Second * 10,
		IdleConnTimeout:       time.Second * 10,
	}

	// get comic page 
	// data, err := client.Comics(3)

	// get comic page by tag
	// data, err := client.ComicsByTagName("group", 1)

	// get comic info
	// data, err := client.ComicInfo(382384)

	// get tags page
	// data, err := client.Tags(1)

	// print comic page picture url
	// println(client.PageUrl(2075216, 2, "j"))
}

```