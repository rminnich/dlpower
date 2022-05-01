// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// from: github.com/u-root/cpu/client

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	config "github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

// UserKeyConfig sets up authentication for a User Key.
// It is required in almost all cases.
func (c *PDU) UserKeyConfig() error {
	kf := c.PrivateKeyFile
	if len(kf) == 0 {
		kf = config.Get(c.Host, "IdentityFile")
		V("key file from config is %q", kf)
		if len(kf) == 0 {
			kf = "key.pub"
		}
	}
	// The kf will always be non-zero at this point.
	if strings.HasPrefix(kf, "~/") {
		kf = filepath.Join(os.Getenv("HOME"), kf[1:])
	}
	key, err := ioutil.ReadFile(kf)
	if err != nil {
		return fmt.Errorf("unable to read private key %q: %v", kf, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("ParsePrivateKey %q: %v", kf, err)
	}
	c.config.Auth = append(c.config.Auth, ssh.PublicKeys(signer))
	return nil
}

// HostKeyConfig sets the host key. It is optional.
func (c *PDU) HostKeyConfig(hostKeyFile string) error {
	hk, err := ioutil.ReadFile(hostKeyFile)
	if err != nil {
		return fmt.Errorf("unable to read host key %v: %v", hostKeyFile, err)
	}
	pk, err := ssh.ParsePublicKey(hk)
	if err != nil {
		return fmt.Errorf("host key %v: %v", string(hk), err)
	}
	c.config.HostKeyCallback = ssh.FixedHostKey(pk)
	return nil
}

// SetEnv sets zero or more environment variables for a Session.
func (c *PDU) SetEnv(envs ...string) error {
	for _, v := range append(os.Environ(), envs...) {
		env := strings.SplitN(v, "=", 2)
		if len(env) == 1 {
			env = append(env, "")
		}
		if err := c.session.Setenv(env[0], env[1]); err != nil {
			return fmt.Errorf("Warning: c.session.Setenv(%q, %q): %v", v, os.Getenv(v), err)
		}
	}
	return nil
}

// Dial implements ssh.Dial for cpu.
// Additionaly, if PDU.Root is not "", it
// starts up a server for 9p requests.
func (c *PDU) Dial() error {
	if err := c.UserKeyConfig(); err != nil {
		return err
	}
	addr := net.JoinHostPort(c.HostName, c.Port)
	cl, err := ssh.Dial(c.network, addr, &c.config)
	V("cpu:ssh.Dial(%s, %s, %v): (%v, %v)", c.network, addr, c.config, cl, err)
	if err != nil {
		return fmt.Errorf("Failed to dial: %v", err)
	}
	c.client = cl
	return nil
}

// GetHostName reads the host name from the ssh config file,
// if needed. If it is not found, the host name is returned.
func GetHostName(host string) string {
	h := config.Get(host, "HostName")
	if len(h) != 0 {
		host = h
	}
	return host
}

// Command creates a PDU
func Command(host string) *PDU {
	return &PDU{
		Host:     host,
		HostName: GetHostName(host),
		Port:     "22",
		network:  "tcp",
		config: ssh.ClientConfig{
			User:            "admin",
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}
