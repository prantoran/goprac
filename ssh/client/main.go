package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SSHAgent() ssh.AuthMethod {

	// first add certificate to ssh-agent
	// ssh-add /path/to/your/private/certificate/file

	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func execute(client *ssh.Client, cmd string) error {
	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("Failed to create session: %s", err)
		os.Exit(1)
	}
	defer session.Close()

	// accessing xterm as pty

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return err
	}

	// setting up pipes between the local process and remote process

	// stdin, err := session.StdinPipe()
	// if err != nil {
	// 	return err

	// }
	// go io.Copy(stdin, os.Stdin)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	fmt.Println("setting env")

	// setting up env
	if err := session.Setenv("LC_USR_DIR", "/usr"); err != nil {
		return err
	}

	fmt.Println("pass stp env")
	err = session.Run(cmd)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	sshConfig := &ssh.ClientConfig{
		User:            "pinku",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			// ssh.Password("your_password"),
			// PublicKeyFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub"),
			SSHAgent(),
		},
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:2222", sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: %s\n", err)
		os.Exit(1)
	}
	defer client.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		fmt.Println("cmd:", cmd)
		if err := execute(client, cmd); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	}

	if scanner.Err() != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	oschan := make(chan os.Signal)

	<-oschan

}
