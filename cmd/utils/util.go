package main

/*
#include <stdint.h>

// Function prototype for the assembly function
extern uint64_t countTrailingZeros(uint64_t x);
*/
import "C"
import (
	"fmt"
)

// Function to count trailing zeros using assembly
func CountTrailingZeros(x uint64) uint64 {
	return uint64(C.countTrailingZeros(C.uint64_t(x)))
}

func main() {
	number := uint64(16) // Example number
	result := CountTrailingZeros(number)
	fmt.Printf("Number of trailing zeros in %d: %d\n", number, result)
}
