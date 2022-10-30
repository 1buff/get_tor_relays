package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/imroc/req/v3"
	"github.com/simonnilsson/ask"
)

func check_bridge(bridge, fingerprint string, bridge_chan chan string) {
	conn, err := net.DialTimeout("tcp", bridge, 500*time.Millisecond)
	if err == nil {
		bridge_chan <- fmt.Sprintf("Bridge %s %s", bridge, fingerprint)
		conn.Close()
	}

}

func get_bridges_and_check() {
	var response map[string]interface{}
	var working_bridges = []string{}
	bridges_chan := make(chan string)

	body, err := req.Get("https://onionoo.torproject.org/details?type=relay&running=true&fields=fingerprint,or_addresses")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body.Bytes(), &response)

	go func() {
		for i := 0; ; i++ {
			select {
			case working_bridge := <-bridges_chan:
				if i >= 5 {
					for _, j := range working_bridges {
						fmt.Println(j)
					}
					fmt.Println("UseBridges 1")
					os.Exit(0)
				}
				working_bridges = append(working_bridges, working_bridge)
			}
		}
	}()

	for i := 0; ; i++ {
		address := ask.For(response, fmt.Sprintf("relays[%d].or_addresses[0]", i)).Value()
		fingerprint := ask.For(response, fmt.Sprintf("relays[%d].fingerprint", i)).Value()

		go check_bridge(address.(string), fingerprint.(string), bridges_chan)

		time.Sleep(50 * time.Millisecond)
	}

}
func main() {
	get_bridges_and_check()
}
