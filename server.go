// server
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
)

//port; port connection string. example ":8080"
func Receive(key []byte, port string) ([]byte, error) {
	cur, err := receiveCourier(port)
	if err != nil {
		fmt.Println("Connection error", err)
		return nil, err
	}

	plainbytes, err := shortdecrypt(key, cur.Nonce, cur.Cipherbytes)
	if err != nil {
		fmt.Println("Decryption error", err)
		return nil, err
	}

	if cur.Compressed {
		plainbytes = shortdecompress(plainbytes)
	}
	return plainbytes, nil
}

//Returns the lzma-decompressed bytes
func shortdecompress(shrunk []byte) []byte {
	b := bytes.NewBuffer(shrunk)
	r := lzma.NewReader(b)
	var bb bytes.Buffer
	io.Copy(&bb, r)
	r.Close()
	return bb.Bytes()
}

//Decrypts (AES-GCM)
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

//Returns a Courier from the network
//port string; for example ":8080"
func receiveCourier(port string) (Courier, error) {
	ln, err := net.Listen("tcp", port)
	p := &Courier{}
	if err != nil {
		return *p, err
	}
	conn, err := ln.Accept() // this blocks until connection or error
	if err != nil {
		return *p, err
	}
	dec := gob.NewDecoder(conn)
	dec.Decode(p)
	conn.Close()
	ln.Close()
	return *p, nil
}
