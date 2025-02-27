//go:build !windows
// +build !windows

package playwright_test

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/idaunis/playwright-go"
)

type remoteServer struct {
	url string
	cmd *exec.Cmd
}

func newRemoteServer() (*remoteServer, error) {
	driver, err := playwright.NewDriver(&playwright.RunOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not start Playwright: %v", err)
	}
	cmd := exec.Command(driver.DriverBinaryLocation, "launch-server", browserName)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get stdout pipe: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start server: %v", err)
	}
	scanner := bufio.NewReader(stdout)
	url, err := scanner.ReadString('\n')
	url = strings.TrimRight(url, "\n")
	if err != nil {
		return nil, fmt.Errorf("could not read url: %v", err)
	}
	return &remoteServer{
		url: url,
		cmd: cmd,
	}, nil
}

func (s *remoteServer) Close() {
	_ = syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL)
	<-time.After(time.Second)
}
