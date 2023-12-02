package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
)

var (
	modeFlag = flag.String("mode", "passive", "Specify the mode (active/passive)")
	url = flag.String("url", "none", "Set the target server for active mode")
	port = flag.String("port", "7980", "Set the target port in passive mode")
)

func main() {
	fmt.Println("****Minecraft banned-players.json Poster****")
	fmt.Println("****Developed by YuxiangWang_0525****\n****Powered by Go Programming Language****")
	
	flag.Parse()
	switch strings.ToLower(*modeFlag) {
	case "active":
		switch strings.ToLower(*url) {
		case "none":
			fmt.Println("Invalid URL")
			os.Exit(1)
		}
		fmt.Println("The program will run in active mode")
		sendBannedPlayers()
	case "passive":
		fmt.Println("The program will run in passive mode")
		startServer()
	default:
		log.Fatal("Invalid parameters")
	}
}

func sendBannedPlayers() {
	for {
		content, err := readBannedPlayersFile()
		if err != nil {
			log.Println("Failed to read the banner players. json file:", err)
			return
		}

		encodedContent := base64.StdEncoding.EncodeToString(content)
		err = postData(encodedContent)
		if err != nil {
			log.Println("Failed to send POST request to server:", err)
		} else {
			fmt.Println("Banned players data has been sent to the server")
		}

		time.Sleep(2 * time.Minute)
	}
}

func readBannedPlayersFile() ([]byte, error) {
	content, err := ioutil.ReadFile("banned-players.json")
	if err != nil {
		return nil, err
	}
	return content, nil
}

func postData(content string) error {
	payload := strings.NewReader("content=" + content)

	req, err := http.NewRequest("POST", *url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			content, err := readBannedPlayersFile()
			if err != nil {
				log.Println("Failed to read the banner players. json file:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(content)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Listening to local port:"+*port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, nil))
}
