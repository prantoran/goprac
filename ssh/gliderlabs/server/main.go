package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/gliderlabs/ssh"
)

func main() {

	fmt.Println("yoyo")

	log.Fatal(ssh.ListenAndServe(":2222",
		func(s ssh.Session) {

			data := []byte{}

			io.ReadFull(s, data)

			pty, _, ptyallowed := s.Pty()

			command := s.Command()

			fmt.Println(
				"\n\nsession:",
				"\n\nUser:", s.User(),
				"\n\nRemoteAddr Network:", s.RemoteAddr().Network(),
				"\n\nRemoteAddr String:", s.RemoteAddr().String(),
				"\n\nLocalAddr Network:", s.LocalAddr().Network(),
				"\n\nLocalAddr String:", s.LocalAddr().String(),
				"\n\nEnviron:", s.Environ(),
				"\n\nCommand:", command,
				"\n\nPublicKey:", s.PublicKey(),
				"\n\nPermissions CriticalOptions:", s.Permissions().CriticalOptions,
				"\n\nPermissions Extensions:", s.Permissions().Extensions,
				"\n\npty:", pty,
				"\n\nptyallowed:", ptyallowed,
			)

			fmt.Println("data:", string(data))

			var cmd *exec.Cmd

			if len(command) == 1 {
				cmd = exec.Command(command[0])
			} else if len(command) > 1 {
				cmd = exec.Command(command[0], command[1:]...)
			}

			if cmd != nil {

				var stdout, stderr bytes.Buffer

				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				err := cmd.Run()

				// out, err := cmd.CombinedOutput()
				if err != nil {
					log.Fatalf("cmd.Run() failed with %s\n", err)
				}
				// fmt.Printf("combined out:\n%s\n", string(out))

				s.Write(stdout.Bytes())
				s.Write(stderr.Bytes())

			}

		},
		ssh.HostKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa")),
		ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			// ssh public key
			data, _ := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/id_rsa.pub")

			fmt.Println("data:", data)
			allowed, _, _, _, _ := ssh.ParseAuthorizedKey(data)
			return ssh.KeysEqual(key, allowed)
		}),
	)

	// ssh.Handle(func(s ssh.Session) {
	// 	io.WriteString(s, "Hello world\n")
	// })

	// log.Fatal(ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile(os.Getenv("HOME")+"/.ssh/id_rsa")))
}
