#!/bin/bash

# Router Performance Test Suite
# This script runs comprehensive benchmarks for the go-wolf router

echo "🐺 Go-Wolf Router Performance Benchmark Suite"
echo "=============================================="

# Change to the project root
cd "$(dirname "$0")/.."

# Ensure we're in the right directory
if [ ! -f "wolf.go" ]; then
    echo "❌ Error: Must be run from the go-wolf project root directory"
    exit 1
fi


RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' 

echo -e "${BLUE}📊 Running router performance benchmarks...${NC}"
echo

# to run benchmark with nice formatting
run_benchmark() {
    local name="$1"
    local pattern="$2"
    
    echo -e "${YELLOW}🔄 Running: $name${NC}"
    echo "----------------------------------------"
    
    # Run the benchmark and capture output
    go test -bench="$pattern" -benchmem -count=3 -timeout=30m ./test/benchmark/ | grep -E "(Benchmark|PASS|FAIL)"
    
    echo
    echo "----------------------------------------"
    echo
}

# run all benchmark
echo -e "${GREEN}🚀 Basic Router Benchmarks${NC}"
run_benchmark "Basic Routing" "BenchmarkBasicRouting"
run_benchmark "Parameter Routing" "BenchmarkParameterRouting"
run_benchmark "Wildcard Routing" "BenchmarkWildcardRouting"
run_benchmark "Middleware" "BenchmarkMiddleware"
run_benchmark "Route Groups" "BenchmarkRouteGroups"
run_benchmark "Complex Routing" "BenchmarkComplexRouting"

echo -e "${GREEN}🔥 Advanced Performance Tests${NC}"
run_benchmark "Static Routes Scaling" "BenchmarkStaticRoutes"
run_benchmark "Parametric Routes" "BenchmarkParametricRoutes"
run_benchmark "Wildcard Routes" "BenchmarkWildcardRoutes"
run_benchmark "Mixed Routes" "BenchmarkMixedRoutes"
run_benchmark "Route Conflicts" "BenchmarkRouteConflicts"
run_benchmark "HTTP Methods" "BenchmarkHTTPMethods"
run_benchmark "Memory Allocation" "BenchmarkMemoryAllocation"
run_benchmark "Concurrent Routing" "BenchmarkConcurrentRouting"
run_benchmark "Large Path Parameters" "BenchmarkLargePathParameters"
run_benchmark "Deep Nesting" "BenchmarkDeepNesting"

echo -e "${GREEN}⚔️ Comparison Benchmarks${NC}"
run_benchmark "Go-Wolf vs Standard Library" "BenchmarkGoWolfVsStdLib"
run_benchmark "Router Scaling" "BenchmarkRouterScaling"
run_benchmark "Parameter Extraction" "BenchmarkParameterExtraction"
run_benchmark "Middleware Overhead" "BenchmarkMiddlewareOverhead"
run_benchmark "Memory Footprint" "BenchmarkMemoryFootprint"
run_benchmark "Concurrent Load" "BenchmarkConcurrentLoad"
run_benchmark "Longest Path" "BenchmarkLongestPath"
run_benchmark "Worst Case" "BenchmarkWorstCase"

echo -e "${GREEN}✅ Benchmark suite completed!${NC}"
echo
echo -e "${BLUE}📈 Performance Summary:${NC}"
echo "- Check the results above for detailed performance metrics"
echo "- Look for ns/op (nanoseconds per operation) - lower is better"
echo "- Check B/op (bytes per operation) - lower is better"
echo "- Review allocs/op (allocations per operation) - lower is better"
echo
echo -e "${YELLOW}💡 Tips:${NC}"
echo "- Run multiple times for consistent results"
echo "- Compare against other frameworks when available"
echo "- Monitor memory allocations for optimization opportunities"
echo "- Use these benchmarks to track performance regressions"
echo

# generate a report
echo -e "${BLUE}📋 Generating benchmark report...${NC}"
REPORT_FILE="benchmark_report_$(date +%Y%m%d_%H%M%S).txt"

echo "Go-Wolf Router Performance Report" > "$REPORT_FILE"
echo "Generated: $(date)" >> "$REPORT_FILE"
echo "========================================" >> "$REPORT_FILE"
echo >> "$REPORT_FILE"

# run all benchmarks and save to report
go test -bench=. -benchmem -count=1 ./test/benchmark/ >> "$REPORT_FILE" 2>&1

echo -e "${GREEN}📄 Report saved to: $REPORT_FILE${NC}"
echo
echo -e "${GREEN}🎉 All done! Use these benchmarks to track and improve router performance.${NC}"
