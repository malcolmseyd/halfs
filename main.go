package main

import (
	"bytes"
	"fmt"
	"halfs/store"
)

const maxBlobSize = 2958

func main() {
	// blob := []byte("https://seyd.ca")
	// blob := []byte{}
	// for i := 0; i < maxBlobSize; i++ {
	// 	blob = append(blob, byte(i%256))
	// }
	blob := make([]byte, maxBlobSize)
	for i := range blob {
		blob[i] = 0xFF
	}
	// fmt.Println("blob", string(blob))
	ref, err := store.Put(blob)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ref:", ref)
	newBlob, err := store.Get(ref)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("blob", string(blob))
	fmt.Println("equal:", bytes.Equal(blob, newBlob))
}
