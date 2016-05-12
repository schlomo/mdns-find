package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/mdns"
)

func main() {
	service := "_apt_proxy._tcp"
	if len(os.Args) > 1 {
		service = os.Args[1]
	}
	search(service)
}

func search(service string) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			// log.Printf("Found new entry:\n%+v", *entry)
			entry_prettyprint, err := json.MarshalIndent(entry, "", "  ")
			if (err != nil) {
				log.Printf("ERROR: %s", err)
			} else {
				log.Printf("Found new entry:\n%s", entry_prettyprint)
			}

			if entry.AddrV4 != nil {
				proxyURL := fmt.Sprintf("http://%s:%d/", entry.AddrV4, entry.Port)
				var err error
				out := ""
				/*				out, err := exec.Command(
									"debconf-set",
									"mirror/http/proxy",
									proxyURL,
								).CombinedOutput()
								log.Printf("debconf-set err: %s\ndebconf-set out: %s\n", err, out) */
				if err == nil {
					log.Printf("Set proxy to %s", proxyURL)
					os.Exit(0)
				} else {
					log.Printf("ERROR %s from debconf-set: %s", err, out)
				}
			}
		}
	}()

	log.Printf("Searching multicast DNS for %s", service)
	mdnsParams := mdns.DefaultParams(service)
	mdnsParams.Entries = entriesCh
	mdnsParams.Timeout = time.Minute
	mdns.Query(mdnsParams)
	close(entriesCh)
}
