package main

import (
	"net/http"
	"log"
	"time"
	"fmt"
	"flag"
	"bufio"
	"encoding/json"
)

func main() {
	var url string
	var tmpjson interface{}


	flag.StringVar(&url, "url", "", "Hystrix Url")

	flag.Parse()

	if len(url) == 0 {
		log.Panicln("Hystrix Url must be set")
	}

	send := make(chan string, 100)

	go func() {
		for msg := range send {
			fmt.Println(msg)
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

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan(){
			text := scanner.Text()
			if len(text)>0{
				text = text[6:]
				if len(text) > 0 {
					if err = json.Unmarshal([]byte(text), &tmpjson); err == nil {
						send <- text
					}
				}
			}
		}
	}
	log.Println("Exit")
}

