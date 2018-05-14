// Package process is the low-level abstraction for a Tor instance.
//
// The standard use is to create a Creator with NewCreator and the path to the Tor executable. The child package
// 'embedded' can be used if Tor is statically linked in the binary. Most developers will prefer the tor package
// adjacent to this one for a higher level abstraction over the process and control port connection.
package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cretz/bine/util"
)

// Process is the interface implemented by Tor processes.
type Process interface {
	// Start starts the Tor process in the background and returns. It is analagous to os/exec.Cmd.Start.
	Start() error
	// Wait waits for the Tor process to exit and returns error if it was not a successful exit.  It is analagous to
	// os/exec.Cmd.Wait.
	Wait() error
}

// Creator is the interface for process creation.
type Creator interface {
	New(ctx context.Context, args ...string) (Process, error)
}

type exeProcessCreator struct {
	exePath string
}

// NewCreator creates a Creator for external Tor process execution based on the given exe path.
func NewCreator(exePath string) Creator {
	return &exeProcessCreator{exePath}
}

func (e *exeProcessCreator) New(ctx context.Context, args ...string) (Process, error) {
	cmd := exec.CommandContext(ctx, e.exePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}

// ControlPortFromFileContents reads a control port file that is written by Tor when ControlPortWriteToFile is set.
func ControlPortFromFileContents(contents string) (int, error) {
	contents = strings.TrimSpace(contents)
	_, port, ok := util.PartitionString(contents, ':')
	if !ok || !strings.HasPrefix(contents, "PORT=") {
		return 0, fmt.Errorf("Invalid port format: %v", contents)
	}
	return strconv.Atoi(port)
}
