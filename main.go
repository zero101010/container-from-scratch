package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
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

func createRootfs() {
	wgetCmd := exec.Command("wget", "https://dl-cdn.alpinelinux.org/alpine/v3.18/releases/x86_64/alpine-minirootfs-3.18.0-x86_64.tar.gz")

	// Config the output
	// wgetCmd.Stdout = os.Stdout
	wgetCmd.Stderr = os.Stderr

	// Run command wget
	if err := wgetCmd.Run(); err != nil {
		fmt.Println("Error to run wget", err)
		return
	}

	// mkdir
	mkdirCmd := exec.Command("mkdir", "rootfs")

	// config output default
	// mkdirCmd.Stdout = os.Stdout
	mkdirCmd.Stderr = os.Stderr

	// Run mkdir
	if err := mkdirCmd.Run(); err != nil {
		fmt.Println("Error to run mkdir:", err)
		return
	}

	//  tar
	tarCmd := exec.Command("tar", "-xzf", "alpine-minirootfs-3.18.0-x86_64.tar.gz", "-C", "rootfs")

	// Config default output
	// tarCmd.Stdout = os.Stdout
	tarCmd.Stderr = os.Stderr

	// Run tar
	if err := tarCmd.Run(); err != nil {
		fmt.Println("Error to run:", err)
		return
	}

	// rm
	rmCmd := exec.Command("rm", "alpine-minirootfs-3.18.0-x86_64.tar.gz")

	// config standard output
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr

	// Run rm
	if err := rmCmd.Run(); err != nil {
		fmt.Println("Error to run tar:", err)
		return
	}

	fmt.Println("Commands applied sucessfuly")
}

func checkDirExist(path string) bool {
	dirPath := path

	// Check if dir rootfs exist to not create another if alredy exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Dir %s doesn't exist.\n", dirPath)
		return false
	} else {
		fmt.Printf("Dir %s exist.\n", dirPath)
		return true
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

	//Create the $PATH/rootfs to use the alpine struture of file and insert inside
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error to get current Dir:", err)
		return
	}
	fullPath := dir + "/rootfs"
	if checkDirExist(fullPath) == false {
		createRootfs()
	}

	// Change to the new root file system.
	if err := syscall.Chroot(fullPath); err != nil {
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
	// Look up the user and group. Here we'll use "guest" as an example.
	userName := "guest"
	u, err := user.Lookup(userName)
	if err != nil {
		fmt.Println("Error looking up user:", err)
		os.Exit(1)
	}

	// Parse the found UID and GID.
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		fmt.Println("Error parsing UID:", err)
		os.Exit(1)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		fmt.Println("Error parsing GID:", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set the UID and GID for the command.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         uint32(uid),
			Gid:         uint32(gid),
			NoSetGroups: true,
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running the child command:", err)
		os.Exit(1)
	}
}
