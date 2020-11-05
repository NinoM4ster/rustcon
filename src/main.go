package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var (
	addr, pass string
	counter    int64 = 1
	scanner          = bufio.NewScanner(os.Stdin)
)

func init() {
RETRYHOST:
	fmt.Print("Connect to: ")
	fmt.Scanln(&addr)
	if addr == "" {
		fmt.Println("Empty host.")
		goto RETRYHOST
	}
RETRYPASS:
	fmt.Print("Password: ")
	fmt.Scanln(&pass)
	if pass == "" {
		fmt.Println("Empty password.")
		goto RETRYPASS
	}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/" + pass}
	print("Connecting to " + u.String() + "...")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	defer c.Close()
	done := make(chan struct{})
	print("Connected.")

	go func() {
		defer close(done)
		for {
			data := &Data{}
			err := c.ReadJSON(&data)
			if err != nil {
				print(err.Error())
				continue
			}
			switch data.Type {
			case "Generic":
				print(data.Message)
			case "Chat":
				chat := &Chat{}
				err := json.Unmarshal([]byte(data.Message), &chat)
				if err != nil {
					log.Println(err)
					continue
				}
				print(fmt.Sprintf("[%s] %s", chat.Username, chat.Message))
			}
		}
	}()

	for scanner.Scan() {
		print("")
		input := scanner.Text()
		if input == "" {
			continue
		}
		data := &Data{Message: input, Identifier: counter, Type: "Command", Name: "WebRcon"}
		err := c.WriteJSON(data)
		if err != nil {
			print(err.Error())
			continue
		}
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
	if s != "" {
		fmt.Println("\r" + s)
	}
	fmt.Print("\r> ")
}
