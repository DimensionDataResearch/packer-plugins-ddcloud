package helpers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path"
	"strings"
)

// OutputHandler is a function which receives lines of piped output from an external tool as they become available.
type OutputHandler func(string)

// The default output handler (does nothing).
var defaultOutputHandler = func(string) {}

// Tool represents an external tool.
type Tool struct {
	// The tool name.
	Name string

	// ExecutablePath is the path of the tool executable.
	ExecutablePath string

	// WorkDir is the working directory where tool will be run.
	WorkDir string

	// OutputHandler is the function that receives lines of output from ovtool as they become available.
	outputHandler OutputHandler
}

// ForTool creates a new external tool helper.
func ForTool(toolExecutable string, workDir string, outputHandler OutputHandler) (tool *Tool, err error) {
	if toolExecutable == path.Base(toolExecutable) {
		toolExecutable, err = exec.LookPath(toolExecutable)
		if err != nil {
			return
		}
	}

	if outputHandler == nil {
		outputHandler = defaultOutputHandler
	}

	tool = &Tool{
		Name:           path.Base(toolExecutable),
		ExecutablePath: toolExecutable,
		WorkDir:        workDir,
		outputHandler:  outputHandler,
	}

	return tool, nil
}

// Run invokes the tool.
func (tool *Tool) Run(args ...string) (success bool, err error) {
	var (
		toolCommand *exec.Cmd
		stdoutPipe  io.ReadCloser
		stderrPipe  io.ReadCloser
	)
	toolCommand = exec.Command(tool.ExecutablePath, args...)
	toolCommand.Dir = tool.WorkDir

	stdoutPipe, err = toolCommand.StdoutPipe()
	if err != nil {
		return
	}
	defer stdoutPipe.Close()

	stderrPipe, err = toolCommand.StderrPipe()
	if err != nil {
		return
	}
	defer stderrPipe.Close()

	log.Printf("Running tool: '%s' %s",
		tool.ExecutablePath,
		strings.Join(args, " "),
	)

	err = toolCommand.Start()
	if err != nil {
		err = fmt.Errorf("Execute tool: Failed to start: %s", err.Error())

		return
	}

	// Pipe output to the caller.
	tool.scanProcessPipes(stdoutPipe, stderrPipe)

	// Pipes will be auto-closed once process is terminated.
	err = toolCommand.Wait()
	if err != nil {
		err = fmt.Errorf("Execute tool: Did not exit cleanly: %s", err.Error())

		return
	}

	if err != nil {
		err = fmt.Errorf("Execute tool: Failed (%s)", err.Error())

		return
	}

	log.Printf("Execute tool: exited cleanly.")

	success = toolCommand.ProcessState.Success()

	return
}

// Scan STDOUT and STDERR pipes for a process.
//
// Calls the supplied PipeHandler once for each line encountered.
func (tool *Tool) scanProcessPipes(stdioPipe io.ReadCloser, stderrPipe io.ReadCloser) {
	go tool.scanPipe(stdioPipe, "STDOUT")
	go tool.scanPipe(stderrPipe, "STDERR")
}

// Scan a process output pipe, and call the supplied PipeHandler once for each line encountered.
func (tool *Tool) scanPipe(pipe io.ReadCloser, pipeName string) {
	lineScanner := bufio.NewScanner(pipe)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		tool.outputHandler(line)
	}

	scanError := lineScanner.Err()
	if scanError != nil {
		log.Printf("Error scanning tool pipe %s: %s",
			pipeName,
			scanError.Error(),
		)
	}

	pipe.Close()
}
