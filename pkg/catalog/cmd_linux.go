package catalog

import (
	"os/exec"

	"golang.org/x/sys/unix"
)

func pluginCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)

	cmd.SysProcAttr = &unix.SysProcAttr{
		Pdeathsig: unix.SIGKILL,
	}
	return cmd
}
