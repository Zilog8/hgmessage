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
func Receive(key []byte, port string) ([]byte, net.Addr, error) {
	cur, from, err := receiveCourier(port)
	if err != nil {
		fmt.Println("Connection error", err)
		return nil, from, err
	}

	plainbytes, err := shortdecrypt(key, cur.Nonce, cur.Cipherbytes)
	if err != nil {
		fmt.Println("Decryption error", err)
		return nil, from, err
	}

	if cur.Compressed {
		plainbytes = shortdecompress(plainbytes)
	}
	return plainbytes, from, nil
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

//Returns a Courier from the network, and the address where it came from
//port string; for example ":8080"
func receiveCourier(port string) (Courier, net.Addr, error) {
	ln, err := net.Listen("tcp", port)
	p := &Courier{}
	if err != nil {
		return *p, nil, err
	}
	conn, err := ln.Accept() // this blocks until connection or error
	if err != nil {
		return *p, nil, err
	}
	dec := gob.NewDecoder(conn)
	dec.Decode(p)
	remote := conn.RemoteAddr()
	conn.Close()
	ln.Close()
	return *p, remote, nil
}
