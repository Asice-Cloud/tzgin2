# TZGin2 调试器使用指南

## 概述

TZGin2 调试器是一个简化版的交互式调试器，类似于 GDB 或 Delve。它提供了基本的调试功能，包括进程控制、断点管理、内存查看等。

## 快速开始

### 1. 构建项目

```bash
# Linux/macOS
go build -o tzgin2 .

# Windows
go build -o tzgin2.exe .

```

### 2. 启动调试器

```bash
# Windows
tzgin2.exe debug

# Linux/macOS
./tzgin2 debug
```

### 3. 创建测试程序

创建一个简单的 Go 程序用于测试：

```go
// test_program.go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hello, Debugger!")
    
    for i := 0; i < 5; i++ {
        result := fibonacci(i)
        fmt.Printf("fibonacci(%d) = %d\n", i, result)
        time.Sleep(500 * time.Millisecond)
    }
    
    fmt.Println("Program finished")
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

编译测试程序：

```bash
# 包含调试信息
go build -gcflags="-N -l" -o test_program test_program.go
```

### 4. 开始调试

在调试器中：

```bash
(tzdb) launch ./test_program
(tzdb) break main.fibonacci
(tzdb) continue
(tzdb) registers
(tzdb) memory 0x401000 32
(tzdb) stack
(tzdb) step
(tzdb) quit
```

## 命令参考

### 基本命令

- `help, h` - 显示帮助信息
- `quit, q, exit` - 退出调试器

### 程序控制

- `launch <program> [args...]` - 启动程序进行调试
- `continue, c` - 继续执行程序
- `step, s` - 单步执行
- `kill` - 终止当前程序
- `detach` - 从程序分离

### 断点管理

- `break <address|function>` - 设置断点
  - `break main` - 在 main 函数设置断点
  - `break fibonacci` - 在 fibonacci 函数设置断点
  - `break 0x401000` - 在地址 0x401000 设置断点
- `delete <address>` - 删除断点
- `breakpoints, info` - 列出所有断点

### 信息查看

- `registers, regs, r` - 显示寄存器值
- `memory <address> [size]` - 显示内存内容
  - `memory 0x7fff12345678 32` - 显示从地址开始的32字节
- `stack, bt` - 显示堆栈跟踪

## 使用示例

### 完整的调试会话

```bash
$ tzgin2 debug
TZGin2 Debugger v1.0
Type 'help' for available commands

(tzdb) launch ./test_program
Process started with PID: 12345

(tzdb) break fibonacci
Breakpoint set at 0x401100 (simulation)

(tzdb) continue
Process exited

(tzdb) registers
Registers:
  rax: 0x1234567890abcdef
  rbx: 0xfedcba0987654321
  rcx: 0x1111111111111111
  rdx: 0x2222222222222222
  rsi: 0x3333333333333333
  rdi: 0x4444444444444444
  rbp: 0x7fff12345678
  rsp: 0x7fff12345000
  rip: 0x401000

(tzdb) memory 0x401000 32
Memory dump at 0x401000:
00401000: 00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f |................|
00401010: 10 11 12 13 14 15 16 17  18 19 1a 1b 1c 1d 1e 1f |................|

(tzdb) stack
Stack trace:
  #0: 0x401000 main.main
  #1: 0x401100 main.fibonacci
  #2: 0x401110 main.fibonacci
  #3: 0x401120 runtime.main

(tzdb) quit
Goodbye!
```

## 注意事项

### 当前实现的限制

1. **简化实现**: 这是一个教学版本，许多功能是模拟的
2. **平台兼容性**: 主要在 Windows 上测试
3. **调试信息**: 需要编译时包含调试信息才能获得更好的体验

### 真实调试的要求

对于真正的调试器实现，需要：

1. **系统权限**: 
   - Windows: 需要 SeDebugPrivilege 权限
   - Linux: 可能需要 root 权限或设置 ptrace_scope

2. **调试信息编译**:
   ```bash
   go build -gcflags="-N -l" -o program main.go
   ```

3. **符号表**: 确保可执行文件包含符号表

## 进阶使用

### 自定义断点

```bash
# 设置多个断点
(tzdb) break main
(tzdb) break fibonacci
(tzdb) break add

# 查看所有断点
(tzdb) breakpoints
```

### 内存分析

```bash
# 查看不同大小的内存块
(tzdb) memory 0x401000 16
(tzdb) memory 0x7fff12345678 64
```

### 调试技巧

1. **设置入口断点**: 在 main 函数设置断点
2. **逐步执行**: 使用 step 命令单步调试
3. **查看状态**: 结合 registers 和 memory 命令
4. **分析调用栈**: 使用 stack 命令了解程序流程

## 扩展开发

如果要扩展此调试器，可以考虑：

1. **真实的系统调用实现**
2. **更好的符号解析**
3. **反汇编功能**
4. **变量查看**
5. **条件断点**
6. **远程调试支持**

## 故障排除

### 常见问题

1. **无法启动程序**: 检查程序路径和权限
2. **断点无效**: 确保程序包含调试信息
3. **权限错误**: 某些操作可能需要管理员权限

### 调试器调试

如果调试器本身有问题，可以：

1. 检查日志输出
2. 确保 Go 版本兼容
3. 验证系统权限
4. 查看错误消息

这个调试器为学习调试器原理提供了良好的起点，可以在此基础上进行扩展和改进。
