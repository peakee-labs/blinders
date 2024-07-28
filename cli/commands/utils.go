package commands

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/fatih/color"
)

// Execute command and get the last stdout line as the result if not error occur
func Execute(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		color.Red("Can't get stdout pipe: %v", err)
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		color.Red("Can't get stderr pipe: %v", err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		color.Red("Can't start: %v", err)
		return "", err
	}

	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	errChan := make(chan error)

	go readPipe(&stdout, stdoutChan, errChan)
	go readPipe(&stderr, stderrChan, errChan)

	var latestText string
	for {
		select {
		case err := <-errChan:
			if err != nil {
				return "", err
			}
		case line, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
			} else {
				latestText = line
				color.Green("-> %v", line)
			}
		case line, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
			} else {
				color.Red("-> %v", line)
			}
		}
		if stdoutChan == nil && stderrChan == nil {
			break
		}
	}

	close(errChan)

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return latestText, nil
}

func readPipe(pipe *io.ReadCloser, ch chan<- string, errCh chan<- error) {
	scanner := bufio.NewScanner(*pipe)
	for scanner.Scan() {
		ch <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		errCh <- err
	}

	close(ch)
}
