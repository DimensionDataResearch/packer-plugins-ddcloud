package ovftool

import (
	"os/exec"
	"sync"
)

var initOnce sync.Once

// Perform one-time initialisation, if required.
func ensureInitialized() {
	initOnce.Do(initialize)
}

// Initialise ovftool package state.
func initialize() {
	var err error
	ExecutablePath, err = exec.LookPath("ovftool")
	if err != nil {
		panic(err)
	}
}
