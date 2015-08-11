package server

import (
	"fmt"
	"os"
	"reflect"

	"github.com/CenturyLinkCloud/clc-sdk/server"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func connect(ip string, creds server.Credentials) error {
	config := &ssh.ClientConfig{
		User: creds.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(creds.Password),
		},
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		return err
	}

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}

	termWidth, termHeight, err := terminal.GetSize(fd)
	defer terminal.Restore(fd, oldState)
	if err != nil {
		return err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&ssh.ExitError{}) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
