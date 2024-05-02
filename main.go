package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"halfs/store"
)

func main() {
	// blob := []byte("https://seyd.ca")
	// blob := []byte{}
	// for i := 0; i < 12456; i++ {
	// 	blob = append(blob, byte(i%256))
	// }
	blob := make([]byte, 1000000)
	rand.Read(blob)
	// blob := make([]byte, 10000)
	// for i := range blob {
	// 	blob[i] = 'A'
	// }
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
