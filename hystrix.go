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
			if len(text)>6{
				text = text[6:]
				if json.Unmarshal([]byte(text), &tmpjson) == nil {
					fmt.Println(text)
				}
			}
		}
	}
	log.Println("Exit")
}

