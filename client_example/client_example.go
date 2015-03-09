// client_example
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Zilog8/hgmessage"
	"time"
)

func main() {
	fmt.Println("Starting to Send()")
	for compressionlevel := 0; compressionlevel < 10; compressionlevel++ {
		message := YourData{"This is a secret message.", compressionlevel, "Hello World!"}
		hgmessage.Send(message, compressionlevel, []byte("yellow submarine"), "localhost:2018")
		fmt.Println("Sent")
	}

	fmt.Println("Starting to Send<-")
	sendChan, err := hgmessage.SendChannel(3, []byte("yellow submarine"), "localhost:2018")
	if err != nil {
		fmt.Println("Error making channel", err)
		return
	}
	for i := 0; i < 10; i++ {
		message := YourData{"This is a secret message.", i, "Hello World!"}
		sendChan <- message
		fmt.Println("Sent")
	}

	sendChan <- nil

	//Wait till channel empty, or we might cut off transmission
	for len(sendChan) > 0 {
		time.Sleep(2 * time.Second)
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
