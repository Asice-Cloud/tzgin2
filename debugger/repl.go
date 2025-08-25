package debugger

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// REPL 交互式调试器界面
type REPL struct {
	Debugger *Debugger
	Scanner  *bufio.Scanner
}

func NewREPL(debugger *Debugger) *REPL {
	return &REPL{
		Debugger: debugger,
		Scanner:  bufio.NewScanner(os.Stdin),
	}
}

func (r *REPL) Start() {
	fmt.Println("TZGin2 Debugger v1.0")
	fmt.Println("Type 'help' for available commands")

	for {
		fmt.Print("(tzdb) ")

		if !r.Scanner.Scan() {
			break
		}

		line := strings.TrimSpace(r.Scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		if err := r.executeCommand(command, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func (r *REPL) executeCommand(command string, args []string) error {

	switch command {
	case "help", "h":
		r.printHelp()

	case "launch", "l":
		if len(args) == 0 {
			return fmt.Errorf("usage: launch <program> [args...]")
		}
		executable := args[0]
		programArgs := args[1:]

		debugger, err := NewDebugger(executable)
		if err != nil {
			return err
		}

		r.Debugger = debugger
		return r.Debugger.Launch(programArgs)

	case "continue", "c":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		return r.Debugger.Continue()

	case "break", "b":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		if len(args) == 0 {
			return fmt.Errorf("usage: break <address|function>")
		}

		target := args[0]
		var address uint64

		// 解析为地址
		if strings.HasPrefix(target, "0x") {
			addr, err := strconv.ParseUint(target[2:], 16, 64)
			if err != nil {
				return fmt.Errorf("invalid address: %s", target)
			}
			address = addr
		} else {
			// 作为函数名查找
			addr, err := r.Debugger.FindFunction(target)
			if err != nil {
				return err
			}
			address = addr
		}

		if err := r.Debugger.SetBreakpoint(address); err != nil {
			return err
		}

		fmt.Printf("Breakpoint set at 0x%x\n", address)

	case "delete", "d":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		if len(args) == 0 {
			return fmt.Errorf("usage: delete <address>")
		}

		addr, err := strconv.ParseUint(args[0], 0, 64)
		if err != nil {
			return fmt.Errorf("invalid address: %s", args[0])
		}

		if err := r.Debugger.RemoveBreakpoint(addr); err != nil {
			return err
		}

		fmt.Printf("Breakpoint removed at 0x%x\n", addr)

	case "step", "s":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		return r.Debugger.Step()

	case "registers", "regs", "r":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}

		registers, err := r.Debugger.GetRegisters()
		if err != nil {
			return err
		}

		fmt.Println("Registers:")
		for name, value := range registers {
			fmt.Printf("  %s: 0x%016x\n", name, value)
		}

	case "memory", "mem", "x":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		if len(args) < 1 {
			return fmt.Errorf("usage: memory <address> [size]")
		}

		addr, err := strconv.ParseUint(args[0], 0, 64)
		if err != nil {
			return fmt.Errorf("invalid address: %s", args[0])
		}

		size := 16 // 默认读取16字节
		if len(args) > 1 {
			if s, err := strconv.Atoi(args[1]); err == nil {
				size = s
			}
		}

		data, err := r.Debugger.ReadMemory(addr, size)
		if err != nil {
			return err
		}

		r.printMemoryDump(addr, data)

	case "stack", "bt":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}

		stackTrace, err := r.Debugger.GetStackTrace()
		if err != nil {
			return err
		}

		fmt.Println("Stack trace:")
		for i, frame := range stackTrace {
			fmt.Printf("  #%d: %s\n", i, frame)
		}

	case "breakpoints", "info":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}

		if len(r.Debugger.Breakpoints) == 0 {
			fmt.Println("No breakpoints set")
		} else {
			fmt.Println("Breakpoints:")
			for addr, bp := range r.Debugger.Breakpoints {
				status := "enabled"
				if !bp.Enabled {
					status = "disabled"
				}
				fmt.Printf("  0x%x (%s)\n", addr, status)
			}
		}

	case "detach":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}

		if err := r.Debugger.Detach(); err != nil {
			return err
		}

		fmt.Println("Detached from process")
		r.Debugger = nil

	case "kill":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}

		if err := r.Debugger.Kill(); err != nil {
			return err
		}

		fmt.Println("Process killed")
		r.Debugger = nil

	case "quit", "q", "exit":
		if r.Debugger != nil && r.Debugger.IsRunning {
			r.Debugger.Detach()
		}
		fmt.Println("Goodbye!")
		os.Exit(0)
	case "print", "printvar":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		if len(args) == 0 {
			return fmt.Errorf("usage: print <varname> [size]")
		}
		varname := args[0]
		size := 8 // 默认读取8字节
		if len(args) > 1 {
			if s, err := strconv.Atoi(args[1]); err == nil {
				size = s
			}
		}
		addr, err := r.Debugger.FindVariableAddress(varname)
		if err != nil {
			return err
		}
		data, err := r.Debugger.ReadMemory(addr, size)
		if err != nil {
			return err
		}
		fmt.Printf("%s (0x%x): ", varname, addr)
		for _, b := range data {
			fmt.Printf("%02x ", b)
		}
		fmt.Println()
		return nil

	case "set", "setvar":
		if r.Debugger == nil {
			return fmt.Errorf("no program loaded")
		}
		if len(args) < 2 {
			return fmt.Errorf("usage: set <varname> <value>")
		}
		varname := args[0]
		valueStr := args[1]
		addr, err := r.Debugger.FindVariableAddress(varname)
		if err != nil {
			return err
		}
		// 只支持写入8字节整数
		value, err := strconv.ParseUint(valueStr, 0, 64)
		if err != nil {
			return fmt.Errorf("invalid value: %s", valueStr)
		}
		data := make([]byte, 8)
		for i := 0; i < 8; i++ {
			data[i] = byte((value >> (8 * i)) & 0xFF)
		}
		if err := r.Debugger.WriteMemory(addr, data); err != nil {
			return err
		}
		fmt.Printf("Set %s (0x%x) = 0x%x\n", varname, addr, value)
		return nil

	default:
		return fmt.Errorf("unknown command: %s", command)
	}

	return nil
}

// printHelp
func (r *REPL) printHelp() {
	fmt.Println(`Available commands:
  help, h                    - Show this help message
  launch, l <program> [args] - Launch a program for debugging
  continue, c                - Continue execution
  break, b <addr|func>       - Set a breakpoint
  delete, d <addr>           - Remove a breakpoint
  step, s                    - Execute one instruction
  registers, regs, r         - Show register values
  memory, mem, x <addr> [size] - Show memory contents
  stack, bt                  - Show stack trace
  breakpoints, info          - List all breakpoints
  detach                     - Detach from process
  kill                       - Kill the process
  quit, q, exit              - Exit debugger

Examples:
  launch ./myprogram arg1 arg2
  break main
  break 0x401000
  memory 0x7fff12345678 32
  `)
}

// printMemoryDump
func (r *REPL) printMemoryDump(startAddr uint64, data []byte) {
	fmt.Printf("Memory dump at 0x%x:\n", startAddr)

	for i := 0; i < len(data); i += 16 {
		addr := startAddr + uint64(i)
		fmt.Printf("%08x: ", addr)

		// 十六进制字节
		for j := 0; j < 16 && i+j < len(data); j++ {
			if j == 8 {
				fmt.Print(" ")
			}
			fmt.Printf("%02x ", data[i+j])
		}

		// 填充空格
		for j := len(data) - i; j < 16; j++ {
			if j == 8 {
				fmt.Print(" ")
			}
			fmt.Print("   ")
		}

		// 打印 ASCII 表示
		fmt.Print(" |")
		for j := 0; j < 16 && i+j < len(data); j++ {
			b := data[i+j]
			if b >= 32 && b <= 126 {
				fmt.Printf("%c", b)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("|")
	}
}
