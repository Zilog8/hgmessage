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
	Send1()

	fmt.Println("Starting to Send<-")
	Send2()
}

func Send1() {
	for compressionlevel := 0; compressionlevel < 10; compressionlevel++ {
		message := YourData{"Secret message:", compressionlevel, "Hello World!"}
		hgmessage.Send(message, compressionlevel, []byte("yellow submarine"), "localhost:2018")
		fmt.Println("Sent")
	}
}

func Send2() {
	//Make channel
	sendChan, err := hgmessage.SendChannel(3, []byte("yellow submarine"), "localhost:2018")
	if err != nil {
		fmt.Println("Error making channel", err)
		return
	}

	//Send over it
	for i := 0; i < 10; i++ {
		message := YourData{"Secret message:", i, "Hello World!"}
		sendChan <- message
		fmt.Println("Sent")
	}

	//Close channel
	sendChan <- nil

	//Wait till channel empty, or we might cut off transmission since this example code exits when we return
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
