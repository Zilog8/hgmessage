// client
package client

import (
	"bytes"
	"code.google.com/p/lzma"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"github.com/Zilog8/hgmessage"
	"net"
)

// Sends data in a compressed, encrypted and verified manner.
// compressionlevel; <1 or >9 means no compression
// key; any 128-, 192-, or 256-bit key
// connection; for example: "localhost:8080"
func Send(data []byte, compressionlevel int, key []byte, connection string) error {
	compressed := false
	if compressionlevel > 0 && compressionlevel < 10 {
		compressed = true
		data = shortcompress(data, compressionlevel)
	}

	cipherbytes, nonce, err := shortencrypt(key, data)
	if err != nil {
		fmt.Println("Encryption error", err)
		return err
	}

	err = sendCourier(cipherbytes, nonce, compressed, connection)
	if err != nil {
		fmt.Println("Connection error", err)
		return err
	}

	return nil
}

//Returns the lzma-compressed bytes
func shortcompress(expanded []byte, level int) []byte {
	var b bytes.Buffer
	w := lzma.NewWriterSizeLevel(&b, int64(len(expanded)), level)
	w.Write(expanded)
	w.Close()
	return b.Bytes()
}

//Encrypts (AES-GCM) a plainbytes with the provided key (128-, 192-, or 256-bit)
//Returns the resulting cipherbytes and the nonce used or an error
func shortencrypt(key, plainbytes []byte) ([]byte, []byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	cipherbytes := gcm.Seal(nil, nonce, plainbytes, nil)
	return cipherbytes, nonce, nil
}

// connection; for example: "localhost:8080"
func sendCourier(message, nonce []byte, compressed bool, connection string) error {
	conn, err := net.Dial("tcp", connection)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(conn)
	p := &hgmessage.Courier{message, nonce, compressed}
	encoder.Encode(p)
	conn.Close()
	return nil
}
