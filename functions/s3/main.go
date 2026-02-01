package main

import (
	"fmt"
	"unsafe"
)

//
//go:wasmimport gojinn host_s3_put
func host_s3_put(keyPtr, keyLen, valPtr, valLen uint32) uint32

//go:wasmimport gojinn host_s3_get
func host_s3_get(keyPtr, keyLen, outPtr, outMaxLen uint32) uint32

func main() {
	key := "hello.txt"
	val := "Hello from Gojinn WASM!"

	res := host_s3_put(
		ptr(key), uint32(len(key)),
		ptr(val), uint32(len(val)),
	)

	if res != 0 {
		fmt.Println(`{"error": "Failed to write to S3"}`)
		return
	}

	buf := make([]byte, 1024)

	n := host_s3_get(
		ptr(key), uint32(len(key)),
		ptrBytes(buf), uint32(len(buf)),
	)

	if n == 0 {
		fmt.Println(`{"error": "Failed to read from S3 or empty file"}`)
		return
	}

	content := string(buf[:n])

	fmt.Printf(
		`{"body": "S3 Check: I wrote '%s' and read '%s' from the bucket!"}`,
		val,
		content,
	)
}

func ptr(s string) uint32 {
	if len(s) == 0 {
		return 0
	}
	return uint32(uintptr(unsafe.Pointer(unsafe.StringData(s))))
}

func ptrBytes(b []byte) uint32 {
	if len(b) == 0 {
		return 0
	}
	return uint32(uintptr(unsafe.Pointer(&b[0])))
}
