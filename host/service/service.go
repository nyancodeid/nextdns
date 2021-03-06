package service

import (
	"errors"
)

type Service interface {
	Type() string

	Install() error
	Uninstall() error

	Status() (Status, error)
	Start() error
	Stop() error
	Restart() error

	ConfigStorer
}

type Config struct {
	Name        string
	DisplayName string
	Description string
	Arguments   []string
}

type Status int

const (
	StatusUnknown Status = iota
	StatusNotInstalled
	StatusRunning
	StatusStopped
)

var (
	ErrNotSuported      = errors.New("system not supported")
	ErrAlreadyInstalled = errors.New("already installed")
	ErrNoInstalled      = errors.New("not installed")
)

type RunMode int

const RunModeEnv = "SERVICE_RUN_MODE"

const (
	// RunModeNone means the current process is not run as a service.
	RunModeNone RunMode = iota

	// RunModeService specifies that the process is running as a service.
	RunModeService
)
