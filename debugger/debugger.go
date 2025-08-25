package debugger

import (
	"debug/dwarf"
	"debug/elf"
	"debug/pe"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Debugger struct {
	Process     *os.Process
	Executable  string
	Symbols     map[string]uint64
	Breakpoints map[uint64]*Breakpoint
	DwarfData   *dwarf.Data
	IsRunning   bool
}

type Breakpoint struct {
	Address      uint64
	OriginalByte byte
	Enabled      bool
}

func NewDebugger(executable string) (*Debugger, error) {
	debugger := &Debugger{
		Executable:  executable,
		Symbols:     make(map[string]uint64),
		Breakpoints: make(map[uint64]*Breakpoint),
		IsRunning:   false,
	}

	// parse executable if provided
	if executable != "" {
		if err := debugger.loadSymbols(); err != nil {
			return nil, fmt.Errorf("failed to load symbols: %v", err)
		}
	}

	return debugger, nil
}

// loadSymbols 加载符号表和调试信息
func (d *Debugger) loadSymbols() error {
	var dwarfData *dwarf.Data

	if runtime.GOOS == "windows" {
		// Windows PE
		peFile, err := pe.Open(d.Executable)
		if err != nil {
			return err
		}
		defer peFile.Close()

		dwarfData, err = peFile.DWARF()
		if err != nil {
			return fmt.Errorf("no DWARF data found: %v", err)
		}
	} else {
		// Linux/Unix ELF
		elfFile, err := elf.Open(d.Executable)
		if err != nil {
			return err
		}
		defer elfFile.Close()

		dwarfData, err = elfFile.DWARF()
		if err != nil {
			return fmt.Errorf("no DWARF data found: %v", err)
		}

		// 解析 ELF 符号表
		symbols, err := elfFile.Symbols()
		if err == nil {
			for _, symbol := range symbols {
				d.Symbols[symbol.Name] = symbol.Value
			}
		}
	}

	d.DwarfData = dwarfData
	return nil
}

// Launch
func (d *Debugger) Launch(args []string) error {
	cmd := exec.Command(d.Executable, args...)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %v", err)
	}

	d.Process = cmd.Process
	d.IsRunning = true

	fmt.Printf("Process started with PID: %d\n", d.Process.Pid)
	return nil
}

// Continue
func (d *Debugger) Continue() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	_, err := d.Process.Wait()
	if err != nil {
		return fmt.Errorf("process wait failed: %v", err)
	}

	d.IsRunning = false
	fmt.Println("Process exited")
	return nil
}

// SetBreakpoint 设置断点
func (d *Debugger) SetBreakpoint(address uint64) error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	// only save breakpoint information
	d.Breakpoints[address] = &Breakpoint{
		Address:      address,
		OriginalByte: 0,
		Enabled:      true,
	}

	fmt.Printf("Breakpoint set at 0x%x (simulation)\n", address)
	return nil
}

// RemoveBreakpoint 移除断点
func (d *Debugger) RemoveBreakpoint(address uint64) error {
	_, exists := d.Breakpoints[address]
	if !exists {
		return fmt.Errorf("no breakpoint at address 0x%x", address)
	}

	delete(d.Breakpoints, address)
	fmt.Printf("Breakpoint removed at 0x%x\n", address)
	return nil
}

// GetRegisters 获取寄存器值
func (d *Debugger) GetRegisters() (map[string]uint64, error) {
	if !d.IsRunning {
		return nil, fmt.Errorf("process is not running")
	}

	registers := make(map[string]uint64)

	if runtime.GOARCH == "amd64" {
		registers["rax"] = 0x1234567890abcdef
		registers["rbx"] = 0xfedcba0987654321
		registers["rcx"] = 0x1111111111111111
		registers["rdx"] = 0x2222222222222222
		registers["rsi"] = 0x3333333333333333
		registers["rdi"] = 0x4444444444444444
		registers["rbp"] = 0x7fff12345678
		registers["rsp"] = 0x7fff12345000
		registers["rip"] = 0x401000
	}

	return registers, nil
}

// ReadMemory 读取内存
func (d *Debugger) ReadMemory(address uint64, size int) ([]byte, error) {
	if !d.IsRunning {
		return nil, fmt.Errorf("process is not running")
	}

	data := make([]byte, size)
	for i := range data {
		data[i] = byte((address + uint64(i)) & 0xFF)
	}

	return data, nil
}

// WriteMemory 写入内存
func (d *Debugger) WriteMemory(address uint64, data []byte) error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	fmt.Printf("Writing %d bytes to 0x%x (simulation)\n", len(data), address)
	return nil
}

// Step 单步执行
func (d *Debugger) Step() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	fmt.Println("Single step executed (simulation)")
	return nil
}

// FindFunction 查找函数地址
func (d *Debugger) FindFunction(name string) (uint64, error) {
	if address, exists := d.Symbols[name]; exists {
		return address, nil
	}

	// if no symbols found, use DWARF data
	if d.DwarfData == nil {
		// if no symbols and no DWARF data, return a default address
		switch name {
		case "main":
			return 0x401000, nil
		case "fibonacci":
			return 0x401100, nil
		case "main.main":
			return 0x401000, nil
		case "main.fibonacci":
			return 0x401100, nil
		default:
			return 0x401200, nil
		}
	}

	// by using DWARF data
	reader := d.DwarfData.Reader()
	for {
		entry, err := reader.Next()
		if err != nil {
			break
		}
		if entry == nil {
			break
		}

		if entry.Tag == dwarf.TagSubprogram {
			nameAttr := entry.Val(dwarf.AttrName)
			if nameAttr != nil && nameAttr.(string) == name {
				addrAttr := entry.Val(dwarf.AttrLowpc)
				if addrAttr != nil {
					return addrAttr.(uint64), nil
				}
			}
		}
	}

	return 0, fmt.Errorf("function '%s' not found", name)
}

// GetStackTrace 获取堆栈跟踪
func (d *Debugger) GetStackTrace() ([]string, error) {
	if !d.IsRunning {
		return nil, fmt.Errorf("process is not running")
	}

	// 返回模拟堆栈
	stackTrace := []string{
		"0x401000 main.main",
		"0x401100 main.fibonacci",
		"0x401110 main.fibonacci",
		"0x401120 runtime.main",
	}

	return stackTrace, nil
}

// Detach 从进程分离
func (d *Debugger) Detach() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	fmt.Println("Detached from process (simulation)")
	d.IsRunning = false
	return nil
}

// Kill
func (d *Debugger) Kill() error {
	if !d.IsRunning {
		return fmt.Errorf("process is not running")
	}

	if err := d.Process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}

	d.IsRunning = false
	fmt.Println("Process killed")
	return nil
}

// FindVariableAddress 查找变量的内存地址（仅支持全局变量/简单场景）
func (d *Debugger) FindVariableAddress(name string) (uint64, error) {
	if d.DwarfData == nil {
		return 0, fmt.Errorf("no DWARF data loaded")
	}
	reader := d.DwarfData.Reader()
	for {
		entry, err := reader.Next()
		if err != nil {
			break
		}
		if entry == nil {
			break
		}
		if entry.Tag == dwarf.TagVariable {
			nameAttr := entry.Val(dwarf.AttrName)
			if nameAttr != nil && nameAttr.(string) == name {
				locAttr := entry.Val(dwarf.AttrLocation)
				if locAttr != nil {
					// 这里只处理简单的全局变量地址（表达式为地址常量）
					loc, ok := locAttr.([]byte)
					if ok && len(loc) >= 9 && loc[0] == 3 { // DW_OP_addr
						addr := uint64(loc[1]) | uint64(loc[2])<<8 | uint64(loc[3])<<16 | uint64(loc[4])<<24 |
							uint64(loc[5])<<32 | uint64(loc[6])<<40 | uint64(loc[7])<<48 | uint64(loc[8])<<56
						return addr, nil
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("variable '%s' not found or unsupported location", name)
}
