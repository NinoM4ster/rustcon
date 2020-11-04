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
	print("connecting to " + u.String())
	// log.Printf("connecting to %s", u.String())

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
				print("read: " + err.Error())
				continue
			}
			// fmt.Println("[debug] " + string(message))
			data := &Data{}
			err = json.Unmarshal(message, &data)
			if err != nil {
				print(err.Error())
				continue
			}
			switch data.Type {
			case "Generic":
				print(data.Message)
				// log.Printf("%s\n", data.Message)
			case "Chat":
				chat := &Chat{}
				err := json.Unmarshal([]byte(data.Message), &chat)
				if err != nil {
					log.Println(err)
					continue
				}
				print(fmt.Sprintf("[%s] %s", chat.Username, chat.Message))
				// fmt.Printf("[%s] %s\n", chat.Username, chat.Message)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		input := scanner.Text()
		if input == "" {
			print("")
		}
		data := &Data{Message: input, Identifier: counter, Name: "WebRcon"}
		bytes, err := json.Marshal(data)
		if err != nil {
			print(err.Error())
			continue
		}
		c.WriteMessage(1, bytes)
		// fmt.Print("> ")
	}

	if scanner.Err() != nil {
		print(scanner.Err().Error())
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

func print(s string) {
	fmt.Print("\r ")
	fmt.Println("\r" + s)
	fmt.Print("> ")
}
