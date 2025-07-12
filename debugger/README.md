# readme by ai
# TZ-Gin 调试器

这是一个类似于 GDB 和 Delve 的交互式调试器，专为 Go 程序设计。

## 功能特性

- **进程控制**: 启动、停止、继续执行程序
- **断点管理**: 设置和删除断点
- **内存访问**: 读取和写入程序内存
- **寄存器查看**: 显示 CPU 寄存器状态
- **堆栈跟踪**: 显示调用栈
- **符号解析**: 通过 DWARF 调试信息解析符号
- **跨平台**: 支持 Windows 和 Linux/macOS

## 使用方法

### 启动调试器

```bash
# 启动交互式调试器
tz-gin debug
```

### 基本命令

```bash
# 加载程序
(tzdb) launch ./myprogram arg1 arg2

# 设置断点
(tzdb) break main
(tzdb) break 0x401000

# 继续执行
(tzdb) continue

# 单步执行
(tzdb) step

# 查看寄存器
(tzdb) registers

# 查看内存
(tzdb) memory 0x7fff12345678 32

# 查看堆栈
(tzdb) stack

# 列出断点
(tzdb) breakpoints

# 删除断点
(tzdb) delete 0x401000

# 退出调试器
(tzdb) quit
```

## 架构设计

### 核心组件

1. **调试器引擎** (`debugger.go`)
   - 进程控制和管理
   - 断点设置和管理
   - 内存和寄存器访问
   - 符号表解析

2. **交互界面** (`repl.go`)
   - 命令行解析
   - 用户交互
   - 输出格式化

3. **平台适配** (`windows.go`, `interface.go`)
   - Windows 特定实现
   - 跨平台接口

### 工作原理

1. **进程创建**: 使用 `ptrace` (Linux/macOS) 或 `CreateProcess` (Windows) 创建可调试进程
2. **断点设置**: 在目标地址写入 `int3` 指令 (0xCC)
3. **事件处理**: 监听调试事件，处理断点命中
4. **符号解析**: 通过 DWARF 调试信息解析函数名和地址

## 示例程序

创建一个简单的测试程序：

```go
// test_program.go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("Hello, Debugger!")
    
    for i := 0; i < 10; i++ {
        result := fibonacci(i)
        fmt.Printf("fibonacci(%d) = %d\n", i, result)
        time.Sleep(100 * time.Millisecond)
    }
}

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

编译并调试：

```bash
# 编译程序（包含调试信息）
go build -gcflags="-N -l" -o test_program test_program.go

# 启动调试器
tz-gin debug

# 在调试器中
(tzdb) launch ./test_program
(tzdb) break fibonacci
(tzdb) continue
(tzdb) registers
(tzdb) memory 0x<address> 32
(tzdb) continue
```

## 调试技巧

### 1. 编译时包含调试信息
```bash
go build -gcflags="-N -l" -o myprogram main.go
```

### 2. 设置符号断点
```bash
(tzdb) break main.main
(tzdb) break fibonacci
```

### 3. 查看变量内存
```bash
# 查看栈上的变量
(tzdb) memory $rsp 64
```

### 4. 分析函数调用
```bash
(tzdb) stack
(tzdb) registers
```

## 限制和注意事项

1. **权限要求**: 在某些系统上需要 root 权限
2. **调试信息**: 需要编译时包含调试信息
3. **平台差异**: Windows 和 Linux 实现可能有所不同
4. **复杂性**: 这是一个简化版本，不支持所有 GDB/Delve 功能

## 扩展功能

可以考虑添加的功能：

- 条件断点
- 观察点 (Watchpoints)
- 反汇编显示
- 远程调试
- 多线程调试
- 更好的 Go 特定支持

## 相关资源

- [GDB 文档](https://www.gnu.org/software/gdb/documentation/)
- [Delve 调试器](https://github.com/go-delve/delve)
- [DWARF 调试格式](http://dwarfstd.org/)
- [ptrace 系统调用](https://man7.org/linux/man-pages/man2/ptrace.2.html)
