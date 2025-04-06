package main

/*
#cgo CXXFLAGS: -I.
#cgo LDFLAGS: -L. -lRingHashWrapper
#include "RingHashWrapper.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// Create a ring hash with 3 replicas
	ringHash := C.create_ring_hash(3)
	defer C.destroy_ring_hash(ringHash)

	// Add nodes
	nodes := []string{"Node1", "Node2", "Node3"}
	for _, node := range nodes {
		cNode := C.CString(node)
		C.add_node(ringHash, cNode)
		C.free(unsafe.Pointer(cNode))
	}

	// Get nodes for keys
	keys := []string{"key1", "key2", "key3", "key4"}
	for _, key := range keys {
		cKey := C.CString(key)
		cNode := C.get_node(ringHash, cKey)
		fmt.Printf("Key: %s is mapped to Node: %s\n", key, C.GoString(cNode))
		C.free(unsafe.Pointer(cKey))
	}
}
