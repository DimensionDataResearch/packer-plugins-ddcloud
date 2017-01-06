package ovftool

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// ExecutablePath is the path of the ovftool executable.
var ExecutablePath string

// OutputHandler is a function which receives lines of piped output from ovftool as they become available.
type OutputHandler func(string)

// The default output handler (does nothing).
var defaultOutputHandler = func(string) {}

// Runner executes ovftool.
type Runner struct {
	// ExecutablePath is the path of the ovftool executable.
	ExecutablePath string

	// WorkDir is the working directory where ovftool will be run.
	WorkDir string

	// OutputHandler is the function that receives lines of output from ovtool as they become available.
	outputHandler OutputHandler
}

// NewRunner creates a new Runner.
func NewRunner(workDir string) *Runner {
	return NewRunnerWithOutputHandler(workDir, nil)
}

// NewRunnerWithOutputHandler creates a new Runner with an OutputHandler.
func NewRunnerWithOutputHandler(workDir string, outputHandler OutputHandler) *Runner {
	ensureInitialized()

	if outputHandler == nil {
		outputHandler = defaultOutputHandler
	}

	return &Runner{
		ExecutablePath: ExecutablePath,
		WorkDir:        workDir,
		outputHandler:  outputHandler,
	}
}

// Run invokes ovftool.
func (runner *Runner) Run(args ...string) (success bool, err error) {
	var (
		ovftoolCommand *exec.Cmd
		stdoutPipe     io.ReadCloser
		stderrPipe     io.ReadCloser
	)
	ovftoolCommand = exec.Command(ExecutablePath, args...)
	ovftoolCommand.Dir = runner.WorkDir

	stdoutPipe, err = ovftoolCommand.StdoutPipe()
	if err != nil {
		return
	}
	defer stdoutPipe.Close()

	stderrPipe, err = ovftoolCommand.StderrPipe()
	if err != nil {
		return
	}
	defer stderrPipe.Close()

	log.Printf("Running ovftool: '%s' %s",
		ExecutablePath,
		strings.Join(args, " "),
	)

	err = ovftoolCommand.Start()
	if err != nil {
		err = fmt.Errorf("Execute ovftool: Failed to start: %s", err.Error())

		return
	}

	// Pipe output to the caller.
	scanProcessPipes(stdoutPipe, stderrPipe, runner.outputHandler)

	// Pipes will be auto-closed once process is terminated.
	err = ovftoolCommand.Wait()
	if err != nil {
		err = fmt.Errorf("Execute ovftool: Did not exit cleanly: %s", err.Error())

		return
	}

	if err != nil {
		err = fmt.Errorf("Execute ovftool: Failed (%s)", err.Error())

		return
	}

	log.Printf("Execute ovftool: exited cleanly.")

	success = ovftoolCommand.ProcessState.Success()

	return
}

// Scan STDOUT and STDERR pipes for a process.
//
// Calls the supplied PipeHandler once for each line encountered.
func scanProcessPipes(stdioPipe io.ReadCloser, stderrPipe io.ReadCloser, pipeOutput OutputHandler) {
	go scanPipe(stdioPipe, pipeOutput, "STDOUT")
	go scanPipe(stderrPipe, pipeOutput, "STDERR")
}

// Scan a process output pipe, and call the supplied PipeHandler once for each line encountered.
func scanPipe(pipe io.ReadCloser, pipeOutput OutputHandler, pipeName string) {
	lineScanner := bufio.NewScanner(pipe)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		pipeOutput(line)
	}

	scanError := lineScanner.Err()
	if scanError != nil {
		log.Printf("Error scanning ovftool pipe %s: %s",
			pipeName,
			scanError.Error(),
		)
	}

	pipe.Close()
}
