//go:build windows
// +build windows

package debugger

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// by AI

// WindowsDebugger Windows 调试器实现
type WindowsDebugger struct {
	Process     *os.Process
	ProcessInfo *syscall.ProcessInformation
	Executable  string
	Breakpoints map[uint64]*Breakpoint
	IsRunning   bool
}

// Windows API 常量
const (
	DEBUG_PROCESS             = 0x00000001
	DEBUG_ONLY_THIS_PROCESS   = 0x00000002
	CREATE_SUSPENDED          = 0x00000004
	EXCEPTION_BREAKPOINT      = 0x80000003
	EXCEPTION_SINGLE_STEP     = 0x80000004
	DBG_CONTINUE              = 0x00010002
	DBG_EXCEPTION_NOT_HANDLED = 0x80010001
	CONTEXT_FULL              = 0x00010007
	THREAD_GET_CONTEXT        = 0x0008
	THREAD_SET_CONTEXT        = 0x0010
	PROCESS_VM_READ           = 0x0010
	PROCESS_VM_WRITE          = 0x0020
	PROCESS_VM_OPERATION      = 0x0008
)

// Windows 系统调用
var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCreateProcess      = kernel32.NewProc("CreateProcessW")
	procWaitForDebugEvent  = kernel32.NewProc("WaitForDebugEvent")
	procContinueDebugEvent = kernel32.NewProc("ContinueDebugEvent")
	procDebugActiveProcess = kernel32.NewProc("DebugActiveProcess")
	procGetThreadContext   = kernel32.NewProc("GetThreadContext")
	procSetThreadContext   = kernel32.NewProc("SetThreadContext")
	procReadProcessMemory  = kernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	procOpenThread         = kernel32.NewProc("OpenThread")
	procCloseHandle        = kernel32.NewProc("CloseHandle")
)

// DEBUG_EVENT 结构体
type DEBUG_EVENT struct {
	DebugEventCode uint32
	ProcessId      uint32
	ThreadId       uint32
	Exception      EXCEPTION_DEBUG_INFO
}

// EXCEPTION_DEBUG_INFO 结构体
type EXCEPTION_DEBUG_INFO struct {
	ExceptionRecord EXCEPTION_RECORD
	FirstChance     uint32
}

// EXCEPTION_RECORD 结构体
type EXCEPTION_RECORD struct {
	ExceptionCode        uint32
	ExceptionFlags       uint32
	ExceptionRecord      uintptr
	ExceptionAddress     uintptr
	NumberParameters     uint32
	ExceptionInformation [15]uintptr
}

// CONTEXT 结构体 (简化版)
type CONTEXT struct {
	ContextFlags                        uint32
	Dr0, Dr1, Dr2, Dr3, Dr6, Dr7        uint64
	FloatSave                           [512]byte
	SegGs, SegFs, SegEs, SegDs          uint32
	Edi, Esi, Ebx, Edx, Ecx, Eax        uint32
	Ebp, Eip, SegCs, EFlags, Esp, SegSs uint32
	ExtendedRegisters                   [512]byte
}

// NewWindowsDebugger 创建 Windows 调试器
func NewWindowsDebugger(executable string) (*WindowsDebugger, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("Windows debugger only works on Windows")
	}

	return &WindowsDebugger{
		Executable:  executable,
		Breakpoints: make(map[uint64]*Breakpoint),
		IsRunning:   false,
	}, nil
}

// Launch 启动程序进行调试
func (d *WindowsDebugger) Launch(args []string) error {
	cmdLine := d.Executable
	for _, arg := range args {
		cmdLine += " " + arg
	}

	var startupInfo syscall.StartupInfo
	var processInfo syscall.ProcessInformation

	startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))

	// 创建进程用于调试
	cmdLinePtr, _ := syscall.UTF16PtrFromString(cmdLine)
	exePtr, _ := syscall.UTF16PtrFromString(d.Executable)

	ret, _, err := procCreateProcess.Call(
		uintptr(unsafe.Pointer(exePtr)),
		uintptr(unsafe.Pointer(cmdLinePtr)),
		0, 0, 0,
		DEBUG_PROCESS|CREATE_SUSPENDED,
		0, 0,
		uintptr(unsafe.Pointer(&startupInfo)),
		uintptr(unsafe.Pointer(&processInfo)),
	)

	if ret == 0 {
		return fmt.Errorf("CreateProcess failed: %v", err)
	}

	d.ProcessInfo = &processInfo
	d.Process = &os.Process{Pid: int(processInfo.ProcessId)}
	d.IsRunning = true

	return nil
}

// Continue 继续执行
func (d *WindowsDebugger) Continue() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	var debugEvent DEBUG_EVENT

	// 等待调试事件
	ret, _, err := procWaitForDebugEvent.Call(
		uintptr(unsafe.Pointer(&debugEvent)),
		syscall.INFINITE,
	)

	if ret == 0 {
		return fmt.Errorf("WaitForDebugEvent failed: %v", err)
	}

	continueStatus := uint32(DBG_CONTINUE)

	// 处理调试事件
	switch debugEvent.DebugEventCode {
	case EXCEPTION_BREAKPOINT:
		fmt.Printf("Breakpoint hit at 0x%x\n", debugEvent.Exception.ExceptionRecord.ExceptionAddress)
		// 这里可以添加断点处理逻辑

	case EXCEPTION_SINGLE_STEP:
		fmt.Println("Single step completed")

	default:
		continueStatus = DBG_EXCEPTION_NOT_HANDLED
	}

	// 继续执行
	procContinueDebugEvent.Call(
		uintptr(debugEvent.ProcessId),
		uintptr(debugEvent.ThreadId),
		uintptr(continueStatus),
	)

	return nil
}

// ReadMemory 读取内存
func (d *WindowsDebugger) ReadMemory(address uint64, size int) ([]byte, error) {
	if !d.IsRunning {
		return nil, fmt.Errorf("process is not running")
	}

	buffer := make([]byte, size)
	var bytesRead uintptr

	ret, _, err := procReadProcessMemory.Call(
		uintptr(d.ProcessInfo.Process),
		uintptr(address),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&bytesRead)),
	)

	if ret == 0 {
		return nil, fmt.Errorf("ReadProcessMemory failed: %v", err)
	}

	return buffer[:bytesRead], nil
}

// WriteMemory 写入内存
func (d *WindowsDebugger) WriteMemory(address uint64, data []byte) error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	var bytesWritten uintptr

	ret, _, err := procWriteProcessMemory.Call(
		uintptr(d.ProcessInfo.Process),
		uintptr(address),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&bytesWritten)),
	)

	if ret == 0 {
		return fmt.Errorf("WriteProcessMemory failed: %v", err)
	}

	return nil
}

// SetBreakpoint 设置断点
func (d *WindowsDebugger) SetBreakpoint(address uint64) error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	// 读取原始字节
	originalByte, err := d.ReadMemory(address, 1)
	if err != nil {
		return fmt.Errorf("failed to read original byte: %v", err)
	}

	// 写入断点指令 (int3 = 0xCC)
	breakpointByte := []byte{0xCC}
	if err := d.WriteMemory(address, breakpointByte); err != nil {
		return fmt.Errorf("failed to write breakpoint: %v", err)
	}

	// 保存断点信息
	d.Breakpoints[address] = &Breakpoint{
		Address:      address,
		OriginalByte: originalByte[0],
		Enabled:      true,
	}

	return nil
}

// RemoveBreakpoint 移除断点
func (d *WindowsDebugger) RemoveBreakpoint(address uint64) error {
	breakpoint, exists := d.Breakpoints[address]
	if !exists {
		return fmt.Errorf("no breakpoint at address 0x%x", address)
	}

	// 恢复原始字节
	originalByte := []byte{breakpoint.OriginalByte}
	if err := d.WriteMemory(address, originalByte); err != nil {
		return fmt.Errorf("failed to restore original byte: %v", err)
	}

	delete(d.Breakpoints, address)
	return nil
}

// GetRegisters 获取寄存器 (简化实现)
func (d *WindowsDebugger) GetRegisters() (map[string]uint64, error) {
	// 这里需要实现获取线程上下文的逻辑
	// 由于复杂性，这里返回一个空的映射
	return make(map[string]uint64), nil
}

// Kill 终止进程
func (d *WindowsDebugger) Kill() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	if err := d.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	d.IsRunning = false
	return nil
}
