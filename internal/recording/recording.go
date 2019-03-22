package recording

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

var exitStatusZero error

// WaitFunc waits for the command execution to stop.
// It returns an error explaining the stop.
type WaitFunc func() error

// StartFunc executes the recording command with the given args and writes the
// output to the provided writer. It returns a function to stop the recording
// and a function to wait for the recording to stop.
type StartFunc func(command string, args []string, w io.Writer) (
	stop context.CancelFunc,
	wait WaitFunc,
)

func run(
	ctx context.Context,
	command string,
	args []string,
	w io.Writer,
	status chan error,
) {
	cmd := exec.CommandContext(ctx, command, args...)
	wait := cmd.Wait
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		status <- err
		close(status)
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		status <- err
		close(status)
		return
	}
	go io.Copy(w, stdout)
	startTime := time.Now()
	log.Println("Recording started")
	err = wait()
	log.Println("Recording stopped")
	canceled := ctx.Err()
	if err != exitStatusZero && canceled != context.Canceled {
		// Command has stopped unexpectedly.
		log.Println(err)
		if time.Since(startTime).Seconds() > 1 {
			// Command ran long enough for this not to be an argument error, restart.
			run(ctx, command, args, w, status)
		} else {
			status <- err
			close(status)
		}
	} else {
		status <- canceled
		close(status)
	}
}

// Start executes the recording command with the given args and writes the
// output to the provided writer. It returns a function to stop the recording
// and a function to wait for the recording to stop.
// If the recording command fails unexpectedly, it is restarted automatically.
func Start(command string, args []string, w io.Writer) (
	stop context.CancelFunc,
	wait WaitFunc,
) {
	ctx, stop := context.WithCancel(context.Background())
	status := make(chan error)
	wait = func() error {
		return <-status
	}
	go run(ctx, command, args, w, status)
	return
}
