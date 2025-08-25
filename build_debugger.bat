@echo off
REM 构建调试器
echo Building tzgin2 debugger...
go build -o tzgin2.exe .

REM 创建测试程序
echo Creating test program...
(
echo package main
echo.
echo import (
echo     "fmt"
echo     "time"
echo ^)
echo.
echo func main^(^) {
echo     fmt.Println^("Hello, Debugger!"^)
echo     
echo     for i := 0; i ^< 5; i++ {
echo         result := fibonacci^(i^)
echo         fmt.Printf^("fibonacci^(%%d^) = %%d\n", i, result^)
echo         time.Sleep^(500 * time.Millisecond^)
echo     }
echo     
echo     fmt.Println^("Program finished"^)
echo }
echo.
echo func fibonacci^(n int^) int {
echo     if n ^<= 1 {
echo         return n
echo     }
echo     return fibonacci^(n-1^) + fibonacci^(n-2^)
echo }
) > test_program.go

REM 编译测试程序（包含调试信息）
echo Compiling test program with debug info...
go build -gcflags="-N -l" -o test_program.exe test_program.go

echo Setup complete!
echo.
echo To start debugging:
echo 1. tzgin2.exe debug
echo 2. In debugger: launch ./test_program.exe
echo 3. In debugger: break fibonacci
echo 4. In debugger: continue
echo.
echo Available commands in debugger:
echo   help - Show all commands
echo   launch ^<program^> - Start program
echo   break ^<function^> - Set breakpoint
echo   continue - Continue execution
echo   step - Single step
echo   registers - Show registers
echo   memory ^<addr^> - Show memory
echo   stack - Show stack trace
echo   quit - Exit debugger

pause
