//go:build !windows
// +build !windows

package debugger

func NewWindowsDebugger(executable string) (DebuggerInterface, error) {
	return nil, nil // stub for non-Windows platforms
}
