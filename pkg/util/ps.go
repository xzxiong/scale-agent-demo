package util

import (
	"fmt"
	"os/exec"
	"strings"
)

type ProcessCmdLine struct {
	Comm string
	Args []string
}

func GetProcessList(inHostNamespace bool) ([]ProcessCmdLine, error) {
	//format := "user,pid,ppid,stime,pcpu,pmem,rss,vsz,stat,time,comm,psr,cgroup"
	format := "comm,cmd"
	out, err := getPsOutput(inHostNamespace, format)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	ret := make([]ProcessCmdLine, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 { // It is the end-line, normally.
			continue
		}
		fmt.Printf("line: %s\n", line)
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("invalid format: %s", line)
		}
		ret = append(ret, ProcessCmdLine{Comm: fields[0], Args: fields[1:]})
	}

	return ret, nil
}

func getPsOutput(inHostNamespace bool, format string) ([]byte, error) {
	args := []string{}
	command := "ps"
	if !inHostNamespace {
		command = "/usr/sbin/chroot"
		args = append(args, "/rootfs", "ps")
	}
	args = append(args, "-e", "-o", format)
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute %q command: %v", command, err)
	}
	return out, err
}
