package main

import (
	"net/http"
	"log"
	"io"
	"encoding/json"
	"time"
	"fmt"
	"flag"
)

func readline(reader io.Reader) (line []byte, err error) {
	line = make([]byte, 0, 100)
	for {
		b := make([]byte, 1)
		n, er := reader.Read(b)
		if n > 0 {
			c := b[0]
			if c == '\n' { // end of line
				break
			}
			line = append(line, c)
		}
		if er != nil {
			err = er
			return
		}
	}
	if len(line) > 6 {
		line = line[6:]
	}
	return
}


func main() {
	var url string
	var tmpjson interface{}


	flag.StringVar(&url, "url", "", "Hystrix Url")

	flag.Parse()

	if len(url) == 0 {
		log.Panicln("Hystrix Url must be set")
	}

	send := make(chan []byte, 100)

	go func() {
		for msg := range send {
			fmt.Println(string(msg))
		}
	}()

	for {
		log.Println("Try get", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Println("Return code is", resp.Status)
			time.Sleep(10 * time.Second)
			continue
		}

		for {
			body, err := readline(resp.Body)
			if err != nil {
				break
			}

			if len(body) == 0 {
				continue
			}

			if err = json.Unmarshal(body, &tmpjson); err == nil {
				send <- body
			}
		}
	}
	log.Println("Exit")
}

