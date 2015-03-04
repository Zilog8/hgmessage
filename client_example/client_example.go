// client_example
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Zilog8/hgmessage"
)

func main() {
	fmt.Println("Starting to send")
	for compressionlevel := 0; compressionlevel < 10; compressionlevel++ {
		plainbytes, err := serializeYourData(YourData{"This is a secret message.", compressionlevel, "Hello World!"})
		if err != nil {
			fmt.Println("Serialization error", err)
			return
		}
		fmt.Println("Ready to send", compressionlevel)
		hgmessage.Send(plainbytes, compressionlevel, []byte("yellow submarine"), "localhost:2018")
		fmt.Println("Sent")
	}
}

type YourData struct {
	A string
	B int
	C string
}

func serializeYourData(y YourData) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(y)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
