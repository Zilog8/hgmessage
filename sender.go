// sender
package hgmessage

import (
	"bytes"
	"code.google.com/p/lzma"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

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

//Blocking Sender:
// compressionlevel; <1 or >9 means no compression
// key; any 128-, 192-, or 256-bit key
// connection; for example: "127.0.0.1:4040"
func Send(message *Box, compressionlevel int, key []byte, connection string) error {
	compressed := false
	messageBytes, _ := message.MarshalBinary()
	if compressionlevel > 0 && compressionlevel < 10 {
		compressed = true
		messageBytes = shortcompress(messageBytes, compressionlevel)
	}

	cipherbytes, nonce, err := shortencrypt(key, messageBytes)
	if err != nil {
		fmt.Println("Encryption error", err)
		return err
	}

	for err = safeSender(cipherbytes, nonce, compressed, connection); err != nil; err = safeSender(cipherbytes, nonce, compressed, connection) {
		fmt.Println("Connection error", err)
	}

	return nil
}

func safeSender(messageBytes, nonce []byte, compressed bool, connection string) error {
	defer func() {
		recover()
	}()

	conn, err := net.Dial("tcp", connection)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(conn)
	p := &courier{Cipherbytes: messageBytes, Nonce: nonce, IsCompressed: compressed}
	encoder.Encode(p)
	conn.Close()
	return nil
}

//Buffered channel version of Send. Pass a nil to close
func SendChannel(compressionlevel int, key []byte, connection string) (chan<- *Box, error) {
	conn, err := net.Dial("tcp", connection)
	if err != nil {
		return nil, err
	}
	gencoder := gob.NewEncoder(conn)

	isCompressed := false
	if compressionlevel > 0 && compressionlevel < 10 {
		isCompressed = true

	}
	c := make(chan *Box, 16)
	go func() {
		for message := <-c; message != nil; message = <-c {
			messageBytes, _ := message.MarshalBinary()
			if isCompressed {
				messageBytes = shortcompress(messageBytes, compressionlevel)
			}
			cipherbytes, nonce, err := shortencrypt(key, messageBytes)
			if err != nil {
				fmt.Println("Encryption error", err)
			}
			p := &courier{Cipherbytes: cipherbytes, Nonce: nonce, IsCompressed: isCompressed}
			for err1 := gencoder.Encode(p); err1 != nil; err1 = gencoder.Encode(p) {
				fmt.Println("Error sending (trying again in 2 seconds):", err1)
				time.Sleep(2 * time.Second)
				if conn != nil {
					conn.Close()
				}
				conn, err1 = net.Dial("tcp", connection)
				if err1 == nil {
					gencoder = gob.NewEncoder(conn)
				}
			}
		}
		conn.Close()
		close(c)
	}()

	return c, nil
}
