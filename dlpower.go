/* SPDX-License-Identifier: GPL-2.0-or-later */

package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"sync"

	"golang.org/x/crypto/ssh"
)

type relay struct {
	Host       string
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

func main() {
	flag.Parse()
	var relays []relay
	var m sync.Mutex
	var wg sync.WaitGroup
	for _, host := range []string{"pdu", "pdu2"} {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			o, err := one(host, "uom get relay/outlets")
			if err != nil {
				log.Fatalf("%v: [%v, %v]", host, o, err)
			}
			var rr []relay
			log.Printf("%v: unmarshall %v", host, string(o))
			if err := json.Unmarshal(o, &rr); err != nil {
				log.Fatalf("unmarshalling %v: %v", string(o), err)
			}
			for i := range rr {
				rr[i].Host = host
			}
			m.Lock()
			relays = append(relays, rr...)
			m.Unlock()
		}(host)
	}
	wg.Wait()
	log.Printf("%d relays %v", len(relays), relays)

}

//func old() {
//	flag.Parse()
//	a := flag.Args()
//	if len(a) < 3 {
//		log.Fatalf("Usage: %q pdu command at-least-one-port", os.Args[0])
//	}
//	pdu := a[0]
//	c, ok := commands[a[1]]
//	if !ok {
//		log.Fatalf("%q: unknown command", a[1])
//	}
//	// Put the ports in the inner loop, commands on the outer loop,
//	// to give things time to set.
//	for _, cmd := range c {
//		for _, port := range a[2:] {
//			stdout, stderr, err := one(pdu, fmt.Sprintf(cmd, port))
//			if err != nil {
//				log.Printf("%q: %q: %q, %v", cmd, port, stderr.String(), err)
//				continue
//			}
//			//fmt.Printf("%q: %q: %q", cmd, port, stdout.String())
//			fmt.Printf("%s", stdout.String())
//		}
//	}
//}
