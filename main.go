package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	cap "github.com/syndtr/gocapability/capability"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	}

}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	fmt.Print(os.Getpid())
	cmd.SysProcAttr = &syscall.SysProcAttr{

		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: 10000, Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: 10000, Size: 1},
		},
	}
	fmt.Println("executing command:", cmd.String(), os.Getpid())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	// fmt.Println("Command executed successfully")
}

func child() {
	process := os.Getpid()
	fmt.Printf("child pid: %d\n", process)

	caps, _ := cap.NewPid2(process)

	// First clear all capability sets
	caps.Clear(cap.EFFECTIVE)
	caps.Clear(cap.AMBIENT)
	caps.Clear(cap.BOUNDING)
	caps.Clear(cap.PERMITTED)

	// Add only the ones you want to all relevant sets
	caps.Set(cap.EFFECTIVE|cap.PERMITTED|cap.BOUNDING,
		cap.CAP_SYS_ADMIN,
		cap.CAP_NET_ADMIN,
		cap.CAP_SYS_CHROOT,
		cap.CAP_SYS_PTRACE,
	)

	// Apply the changes
	// Apply the changes, including the Bounding set
if err := caps.Apply(cap.EFFECTIVE | cap.PERMITTED | cap.BOUNDING); err != nil {
    panic(err)
}

	fmt.Printf("Capabilities applied: %s\n", caps.String())

	// Optional: enter a minimal environment
	if err := syscall.Sethostname([]byte("container")); err != nil {
		panic(err)
	}
	// syscall.Chroot("/home/low/rootfs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	syscall.Unmount("proc", 0)
}


