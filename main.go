package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
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
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	// cmd.Dir = "/home/rtvkiz"
	fmt.Println("executing command:", cmd.String(), os.Getpid())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := syscall.Sethostname([]byte("container")); err != nil {
		panic(err)
	}
	cg()
	syscall.Chroot("/root/Learning/rootfs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	err := cmd.Run()
	syscall.Unmount("proc", 0)
	if err != nil {
		panic(err)
	}

	// fmt.Println("Command executed successfully")

}

// changed the implementation based on CGROUPv2

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "rtvkiz"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "pids.max"), []byte("30"), 0700))
	// Removes the new cgroup in place after the container exits
	// must(ioutil.WriteFile(filepath.Join(pids, "notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
