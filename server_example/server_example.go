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

	letterChannel, err := hgmessage.ReceiveChannel([]byte("yellow submarine"), 2018, "", deserializeYourData)
	if err != nil {
		fmt.Println("Error making channel: ", err)
		return
	}

	for {
		letter := <-letterChannel

		if letter.Data == nil {
			fmt.Println("Deserialization error from", letter.From, ":", err)
			continue
		}
		y := letter.Data.(YourData)
		fmt.Println("Message from", letter.From, ":", y.A, y.B, y.C)
	}
}

type YourData struct {
	A string
	B int
	C string
}

func (y YourData) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(y.A)
	err = enc.Encode(y.B)
	err = enc.Encode(y.C)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func deserializeYourData(b []byte) (hgmessage.Message, error) {
	bb := bytes.NewBuffer(b)
	dec := gob.NewDecoder(bb)
	var y YourData
	err := dec.Decode(&y.A)
	err = dec.Decode(&y.B)
	err = dec.Decode(&y.C)
	if err != nil {
		return nil, err
	}
	return y, nil
}
