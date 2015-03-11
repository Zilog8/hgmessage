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
	"strconv"
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
// port; port to open the channel on
// from; who to accept from. Matches as a string preffix. Example: "127.0." matches "127.0.0.1:50437"
//Returns a channel of Letter, which contains the Data Message and a string interpretation of who sent it.
func ReceiveChannel(key []byte, port int, from string, mum MessageUnmarshaler) (<-chan Letter, error) {

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	courierChannel := make(chan courier, 16)

	go func() {
		for {
			safeListener(ln, from, courierChannel)
		}
	}()

	letterChannel := make(chan Letter, 16)

	go func() {
		for {
			cur := <-courierChannel
			plainbytes, err := shortdecrypt(key, cur.Nonce, cur.Cipherbytes)
			if err != nil {
				fmt.Println("Decryption error:", err)
			}

			if cur.IsCompressed {
				plainbytes = shortdecompress(plainbytes)
			}

			mess, err := mum(plainbytes)
			if err != nil {
				fmt.Println("Unmarshal error:", err)
			} else {
				letter := Letter{Data: mess, From: cur.From}
				letterChannel <- letter
			}
		}
	}()

	return letterChannel, nil
}

func safeListener(ln net.Listener, from string, courierChannel chan<- courier) {
	defer func() {
		recover()
	}()

	conn, err := ln.Accept()
	if err != nil || !strings.HasPrefix(conn.RemoteAddr().String(), from) {
		conn.Close()
	} else {
		dec := gob.NewDecoder(conn)
		p := courier{}
		for err1 := dec.Decode(&p); err1 == nil; err1 = dec.Decode(&p) {
			p.From = conn.RemoteAddr().String()
			courierChannel <- p
			p = courier{}
		}
	}
}
