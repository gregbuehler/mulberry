package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type Mulberry struct {
	Config struct {
		Bind    string `json:"bind"`
		Expires int    `json:"expires"`
	} `json:"config"`
	Groups []struct {
		Dest   []string `json:"dest"`
		Source string   `json:"source"`
	} `json:"groups"`
}

var App Mulberry

func main() {
	fmt.Println("Mulb-errrrry!")

	fmt.Println("Loading ./config.json")
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(file, &App)

	fmt.Printf("Config.Bind: %s\n", App.Config.Bind)
	fmt.Printf("Config.Expires: %d\n", App.Config.Expires)

	for i, g := range App.Groups {
		for j, d := range g.Dest {
			fmt.Printf("[%d] %s => [%d] %s\n", i, g.Source, j, d)
		}
	}

	l, err := net.Listen("tcp", App.Config.Bind)
	if err != nil {
		fmt.Println("Error binding managment: %v\n", err)
		os.Exit(1)
	}

	defer l.Close()

	fmt.Printf("Mulberry managment listening on %s", App.Config.Bind)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting managment connection: %v\n", err)
			os.Exit(1)
		}

		go handleManagementRequest(conn)
	}
}

func handleManagementRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error reading management buffer: %v\n", err)
	}
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}
