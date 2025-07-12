package debugger

import (
	"fmt"
	"runtime"
)

type DebuggerInterface interface {
	Launch(args []string) error
	Continue() error
	SetBreakpoint(address uint64) error
	RemoveBreakpoint(address uint64) error
	ReadMemory(address uint64, size int) ([]byte, error)
	WriteMemory(address uint64, data []byte) error
	GetRegisters() (map[string]uint64, error)
	Kill() error
}

func NewPlatformDebugger(executable string) (DebuggerInterface, error) {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsDebugger(executable)
	case "linux", "darwin":
		return NewDebugger(executable)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
