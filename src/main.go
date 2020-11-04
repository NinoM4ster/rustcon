package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var (
	addr          = flag.String("h", "localhost:28016", "http service address")
	pass          = flag.String("p", "", "RCON password")
	counter int64 = 1
)

func init() {
	flag.Parse()
	log.SetFlags(0)
	if *pass == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/" + *pass}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				continue
			}
			data := &Data{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("[%+v] %+v\n", data.Identifier, data.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		input := scanner.Text()
		data := &Data{Message: input, Identifier: counter, Name: "WebRcon"}
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			continue
		}
		c.WriteMessage(1, bytes)
		// fmt.Print("> ")
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case <-interrupt:
	// 		log.Println("received interrupt")
	// 		// err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		// if err != nil {
	// 		// 	log.Println("write close:", err)
	// 		// 	return
	// 		// }
	// 		// select {
	// 		// case <-done:
	// 		// 	fmt.Println("connection closed cleanly.")
	// 		// case <-time.After(time.Second * 1):
	// 		// 	fmt.Println("waited too long for closing cleanly.")
	// 		// }
	// 		return
	// 	}
	// }

}
