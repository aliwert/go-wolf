@echo off
setlocal enabledelayedexpansion

REM Router Performance Test Suite for Windows
REM This script runs comprehensive benchmarks for the go-wolf router

echo 🐺 Go-Wolf Router Performance Benchmark Suite
echo ==============================================

REM Change to the project root
cd /d "%~dp0\.."

REM Ensure we're in the right directory
if not exist "wolf.go" (
    echo ❌ Error: Must be run from the go-wolf project root directory
    exit /b 1
)

echo 📊 Running router performance benchmarks...
echo.

REM Function to run benchmark (simulated with goto)
goto :start_benchmarks

:run_benchmark
set "name=%~1"
set "pattern=%~2"

echo 🔄 Running: %name%
echo ----------------------------------------

REM Run the benchmark
go test -bench="%pattern%" -benchmem -count=3 -timeout=30m ./test/benchmark/ | findstr /C:"Benchmark" /C:"PASS" /C:"FAIL"

echo.
echo ----------------------------------------
echo.
goto :eof

:start_benchmarks

echo 🚀 Basic Router Benchmarks
call :run_benchmark "Basic Routing" "BenchmarkBasicRouting"
call :run_benchmark "Parameter Routing" "BenchmarkParameterRouting"
call :run_benchmark "Wildcard Routing" "BenchmarkWildcardRouting"
call :run_benchmark "Middleware" "BenchmarkMiddleware"
call :run_benchmark "Route Groups" "BenchmarkRouteGroups"
call :run_benchmark "Complex Routing" "BenchmarkComplexRouting"

echo 🔥 Advanced Performance Tests
call :run_benchmark "Static Routes Scaling" "BenchmarkStaticRoutes"
call :run_benchmark "Parametric Routes" "BenchmarkParametricRoutes"
call :run_benchmark "Wildcard Routes" "BenchmarkWildcardRoutes"
call :run_benchmark "Mixed Routes" "BenchmarkMixedRoutes"
call :run_benchmark "Route Conflicts" "BenchmarkRouteConflicts"
call :run_benchmark "HTTP Methods" "BenchmarkHTTPMethods"
call :run_benchmark "Memory Allocation" "BenchmarkMemoryAllocation"
call :run_benchmark "Concurrent Routing" "BenchmarkConcurrentRouting"
call :run_benchmark "Large Path Parameters" "BenchmarkLargePathParameters"
call :run_benchmark "Deep Nesting" "BenchmarkDeepNesting"

echo ⚔️ Comparison Benchmarks
call :run_benchmark "Go-Wolf vs Standard Library" "BenchmarkGoWolfVsStdLib"
call :run_benchmark "Router Scaling" "BenchmarkRouterScaling"
call :run_benchmark "Parameter Extraction" "BenchmarkParameterExtraction"
call :run_benchmark "Middleware Overhead" "BenchmarkMiddlewareOverhead"
call :run_benchmark "Memory Footprint" "BenchmarkMemoryFootprint"
call :run_benchmark "Concurrent Load" "BenchmarkConcurrentLoad"
call :run_benchmark "Longest Path" "BenchmarkLongestPath"
call :run_benchmark "Worst Case" "BenchmarkWorstCase"

echo ✅ Benchmark suite completed!
echo.
echo 📈 Performance Summary:
echo - Check the results above for detailed performance metrics
echo - Look for ns/op (nanoseconds per operation) - lower is better
echo - Check B/op (bytes per operation) - lower is better
echo - Review allocs/op (allocations per operation) - lower is better
echo.
echo 💡 Tips:
echo - Run multiple times for consistent results
echo - Compare against other frameworks when available
echo - Monitor memory allocations for optimization opportunities
echo - Use these benchmarks to track performance regressions
echo.

REM Generate a simple report
echo 📋 Generating benchmark report...

REM Create timestamp for report file
for /f "tokens=2 delims==" %%a in ('wmic OS Get localdatetime /value') do set "dt=%%a"
set "YY=%dt:~2,2%" & set "YYYY=%dt:~0,4%" & set "MM=%dt:~4,2%" & set "DD=%dt:~6,2%"
set "HH=%dt:~8,2%" & set "Min=%dt:~10,2%" & set "Sec=%dt:~12,2%"
set "timestamp=%YYYY%%MM%%DD%_%HH%%Min%%Sec%"

set "REPORT_FILE=benchmark_report_%timestamp%.txt"

echo Go-Wolf Router Performance Report > "%REPORT_FILE%"
echo Generated: %date% %time% >> "%REPORT_FILE%"
echo ======================================== >> "%REPORT_FILE%"
echo. >> "%REPORT_FILE%"

REM Run all benchmarks and save to report
go test -bench=. -benchmem -count=1 ./test/benchmark/ >> "%REPORT_FILE%" 2>&1

echo 📄 Report saved to: %REPORT_FILE%
echo.
echo 🎉 All done! Use these benchmarks to track and improve router performance.

pause
