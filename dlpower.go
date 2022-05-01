/* SPDX-License-Identifier: GPL-2.0-or-later */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type relay struct {
	Host       string
	number     int
	Name       string `json:"name"`
	Tstate     bool   `json:"transient_state"`
	Critical   bool   `json:"critical"`
	Pstate     bool   `json:"physical_state"`
	Locked     bool   `json:"locked"`
	State      bool   `json:"state"`
	CycleDelay string `json:"cycle_delay"`
}

// PDU defines a PDU.
type PDU struct {
	config  ssh.ClientConfig
	client  *ssh.Client
	session *ssh.Session
	Stdout  io.Reader
	Stderr  io.Reader
	Host    string
	// HostName as found in .ssh/config; set to Host if not found
	HostName       string
	Args           []string
	HostKeyFile    string
	PrivateKeyFile string
	Port           string
	Env            []string
	network        string // This is a variable but we expect it will always be tcp
	cmd            string // The command is built up, bit by bit, as we configure the client
}

var (
	// V is the printer for debug messages.
	V = func(string, ...interface{}) {}
)

func one(host, cmd string) ([]byte, error) {
	c := Command(host)
	V("c %v", c)
	c.cmd = cmd
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	return session.CombinedOutput(cmd)
}

func getRelays(hosts ...string) ([]relay, error) {
	var relays []relay
	var m sync.Mutex
	var wg sync.WaitGroup
	echan := make(chan error)
	for _, host := range hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			o, err := one(host, "uom get relay/outlets")
			if err != nil {
				echan <- fmt.Errorf("%v: [%v, %s]", host, o, err)
				return
			}
			var rr []relay
			V("%v: unmarshall %v", host, string(o))
			if err := json.Unmarshal(o, &rr); err != nil {
				echan <- fmt.Errorf("unmarshalling %v: %w", string(o), err)
				return
			}
			for i := range rr {
				rr[i].number = i
				rr[i].Host = host
			}
			m.Lock()
			relays = append(relays, rr...)
			m.Unlock()
		}(host)
	}
	wg.Wait()
	close(echan)
	var errcount int
	for i := range echan {
		log.Printf("err %v", i)
		errcount++
	}
	return relays, nil
}
func main() {
	var (
		debug  = flag.Bool("d", false, "enable debugging")
		dryrun = flag.Bool("dryrun", false, "dryrun mode")
	)
	flag.Parse()

	if *debug {
		V = log.Printf
	}
	relays, err := getRelays("pdu", "pdu2")
	if err != nil {
		log.Fatal(err)
	}
	V("%d relays %v", len(relays), relays)
	a := flag.Args()
	if len(a) == 0 {
		for _, r := range relays {
			fmt.Printf("Host: %v, Name: %v\n", r.Host, r.Name)
		}
		return
	}
	// If the argv is just a command, the pattern is .
	pat := "."
	if len(a) > 1 {
		pat = strings.Join(a[1:], "|")
	}
	relay := regexp.MustCompile(pat)

	cmd := a[0]
	var printer func(relay int) string
	switch cmd {
	case "on":
		printer = func(relay int) string {
			return fmt.Sprintf("uom set relay/outlets/%d/transient_state true", relay)
		}
	case "off":
		printer = func(relay int) string {
			return fmt.Sprintf("uom set relay/outlets/%d/transient_state false", relay)
		}
	case "cycle":
		printer = func(relay int) string {
			return fmt.Sprintf("uom invoke relay/outlets/%d/cycle", relay)
		}
	default:
		log.Fatalf("%s is not a valid command: use on of on, off, cycle", cmd)
	}

	for _, r := range relays {
		if !relay.MatchString(r.Name) {
			continue
		}
		if *dryrun {
			log.Printf("[dryrun]: %v: relay %v(%d): %v", r.Host, r.Name, r.number, printer(r.number))
			continue
		}
		o, err := one(r.Host, printer(r.number))
		if err != nil {
			log.Printf("%v: relay %v(%d): %v, %v", r.Host, r.Name, r.number, string(o), err)
		}
		if *debug {
			log.Printf("%v: relay %v(%d): %v", r.Host, r.Name, r.number, string(o))
		}
	}

}
