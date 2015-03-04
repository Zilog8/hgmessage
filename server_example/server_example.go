// server_example
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Zilog8/hgmessage"
)

func main() {
	fmt.Println("Server start")
	for {
		serial, err := hgmessage.Receive([]byte("yellow submarine"), ":2018")
		if err != nil {
			fmt.Println("Reception error", err)
			continue
		}

		y, err := deserializeYourData(serial)
		if err != nil {
			fmt.Println("Deserialization error", err)
			continue
		}

		fmt.Println("Decoded to :", y.A, y.B, y.C)
	}
}

type YourData struct {
	A string
	B int
	C string
}

func deserializeYourData(b []byte) (YourData, error) {
	bb := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bb)
	var y YourData
	err := dec.Decode(&y)
	if err != nil {
		return y, err
	}
	return y, nil
}
