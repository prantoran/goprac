package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker         run image <cmd> <params>
// go run main.go run       <cmd> <params>
func main() {

	if err := syscall.Mount("none", flag.Args()[0], "proc", 0, ""); err != nil {
		fmt.Fprintf(os.Stderr, "unshare: mount %v failed: %v", os.Args, err)
		os.Exit(2)
	}

	switch os.Args[1] { // go run main.go [os.Args[1]]
	case "run":
		fmt.Println("run called")
		run()
	case "child":
		fmt.Println("child called")
		child()
	default:
		panic("bad command")
	}
}

func run() {

	// create the

	fmt.Printf("Run Running %v\n", os.Args[2:])

	args := append([]string{"child"}, os.Args[2:]...)
	fmt.Println("augmented cmd: ", args)

	// run this program again, run itself
	// run will re-invoke this process inside a new namespace
	// 2nd time child will be called
	cmd := exec.Command("/proc/self/exe", args...)
	// cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// // creating namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// clones the process where we will run our command
		// flag for unix time sharing system namespace, basically our own
		// hostname inside the container, so that the container
		// cannot see what is happening in the host
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	// this process will clone a new process with the new namespace and
	// then create another process in which to execut the args command

	cmd.Run()
}

func child() {
	fmt.Printf("Child Running %v\n", os.Args[2:])

	syscall.Sethostname([]byte("container"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}
