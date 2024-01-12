package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what??")
	}
}

func run() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running the /proc/self/exe command:", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	// Set hostname of the new UTS namespace
	if err := syscall.Sethostname([]byte("IgorContainer")); err != nil {
		fmt.Println("Error setting hostname:", err)
		os.Exit(1)
	}
	// Change to the new root file system.
	if err := syscall.Chroot("/home/ubuntu/rootfs"); err != nil {
		fmt.Println("Error changing root:", err)
		os.Exit(1)
	}
	// Change working directory after changing the root.
	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error changing working directory:", err)
		os.Exit(1)
	}
	// Mount proc. This needs to be done after chroot and chdir.
	// Because I need the /proc inside of our container that was Copy by chroot and chdir
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting proc:", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running the child command:", err)
		os.Exit(1)
	}
}
