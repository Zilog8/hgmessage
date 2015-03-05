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

	serialchan, err := hgmessage.ReceiveChannel([]byte("yellow submarine"), ":2018", "")
	if err != nil {
		fmt.Println("Error making channel", err)
		return
	}

	for {
		box := <-serialchan
		y, err := deserializeYourData(box.Data)
		if err != nil {
			fmt.Println("Deserialization error from", ":", err)
			continue
		}

		fmt.Println("Message from", box.From, ":", y.A, y.B, y.C)
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
