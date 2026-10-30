package main

import (
	"testing"
)



// 1. Benchmark for the original Map implementation
func BenchmarkParseBracketsOriginal(b *testing.B) {
	input := `{"name": "Alice", "age": 30, "skills": ["Go", "Python"], "active": true}`
	b.ResetTimer()
	for b.Loop() {
		parseBrackets(input)
	}
}

// 2. Benchmark for the Optimized Switch implementation
func BenchmarkParseBracketsOptimized(b *testing.B) {
	input := `{"name": "Alice", "age": 30, "skills": ["Go", "Python"], "active": true}`
	b.ResetTimer()
	for b.Loop() {
		parseBracketsOptimized(input)
	}
}


// go test -bench=. -benchmem