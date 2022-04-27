/* SPDX-License-Identifier: GPL-2.0-or-later */

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

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

func one(host, cmd string) (bytes.Buffer, bytes.Buffer, error) {
	c := Command("pdu")
	V("c %v", c)
	c.cmd = cmd
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	if err := c.Start(); err != nil {
		log.Fatal(err)
	}
	b, err := c.Outputs()
	return b[0], b[1], err
}

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) < 2 {
		log.Fatalf("Usage: %q command at-least-one-port", os.Args[0])
	}
	c, ok := commands[a[0]]
	if !ok {
		log.Fatalf("%q: unknown command", a[0])
	}
	// Put the ports in the inner loop, commands on the outer loop,
	// to give things time to set.
	for _, cmd := range c {
		for _, port := range a[1:] {
			stdout, stderr, err := one("pdu", fmt.Sprintf(cmd, port))
			if err != nil {
				log.Printf("%q: %q: %q, %v", cmd, port, stderr.String(), err)
				continue
			}
			//fmt.Printf("%q: %q: %q", cmd, port, stdout.String())
			fmt.Printf("%s", stdout.String())
		}
	}
}
