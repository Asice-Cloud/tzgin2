#!/bin/bash
echo "[INFO] Building tzgin2 debugger (仅支持 Linux/amd64)..."
GOOS=linux GOARCH=amd64 go build -o tzgin2 .
echo "Building tzgin2 debugger..."
go build -o tzgin2 .

# 创建测试程序
echo "Creating test program..."
cat > test_program.go << 'EOF'
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
EOF

# 编译测试程序（包含调试信息）
echo "Compiling test program with debug info..."
go build -gcflags="-N -l" -o test_program test_program.go

echo "Setup complete!"
echo ""
echo "To start debugging:"
echo "1. ./tzgin2 debug"
echo "2. In debugger: launch ./test_program"
echo "3. In debugger: break main.fibonacci"
echo "4. In debugger: continue"
echo ""
echo "Available commands in debugger:"
echo "  help - Show all commands"
echo "  launch <program> - Start program"
echo "  break <function> - Set breakpoint"
echo "  continue - Continue execution"
echo "  step - Single step"
echo "  registers - Show registers"
echo "  memory <addr> - Show memory"
echo "  stack - Show stack trace"
echo "  quit - Exit debugger"
