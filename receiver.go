// receiver
package hgmessage

import (
	"bytes"
	"code.google.com/p/lzma"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strings"
)

//Returns the lzma-decompressed bytes
func shortdecompress(shrunk []byte) []byte {
	b := bytes.NewBuffer(shrunk)
	r := lzma.NewReader(b)
	var bb bytes.Buffer
	io.Copy(&bb, r)
	r.Close()
	return bb.Bytes()
}

//Decrypts and Authenticates (AES-GCM)
//Returns the resulting plainbytes or an error
func shortdecrypt(key, nonce, cipherbytes []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return nil, err
	}

	plainbytes, err := gcm.Open(nil, nonce, cipherbytes, nil)
	if err != nil {
		return nil, err
	}

	return plainbytes, nil
}

//Receiver Channel:
// key; the 128-, 192-, or 256-bit key used to encrypt
// port; for example: ":4040"
// from; who to accept from. Matches as a string preffix. Example: "127.0." matches "127.0.0.1:50437"
//Returns a channel of *Box, which contains the data []byte and a string interpretation of who sent it.
func ReceiveChannel(key []byte, port string, from string) (<-chan *Box, error) {

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	courierchan := make(chan *courier, 16)

	go func() {
		for {
			safeListener(ln, from, courierchan)
		}
	}()

	boxchan := make(chan *Box, 16)

	go func() {
		for {
			cur := <-courierchan
			plainbytes, err := shortdecrypt(key, cur.Nonce, cur.Cipherbytes)
			if err != nil {
				fmt.Println("Decryption error", err)
			}

			if cur.IsCompressed {
				plainbytes = shortdecompress(plainbytes)
			}
			boxchan <- &Box{Data: plainbytes, From: cur.From}
		}
	}()

	return boxchan, nil
}

func safeListener(ln net.Listener, from string, courchan chan<- *courier) {
	defer func() {
		recover()
	}()

	conn, err := ln.Accept() // this blocks until connection or error
	if err != nil || !strings.HasPrefix(conn.RemoteAddr().String(), from) {
		conn.Close()
	} else {
		dec := gob.NewDecoder(conn)
		p := &courier{}
		for err1 := dec.Decode(p); err1 == nil; err1 = dec.Decode(p) {
			p.From = conn.RemoteAddr().String()
			courchan <- p
			p = &courier{}
		}
	}
}
