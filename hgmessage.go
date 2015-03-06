package hgmessage

import (
	"bytes"
	"encoding/binary"
)

type courier struct {
	Cipherbytes, Nonce []byte
	IsCompressed       bool

	//Not serialized; for use by receiver only
	From               string
}

func (c *courier) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	//IsCompressed
	isByte := 0
	if c.IsCompressed {
		isByte = 1
	}
	b.WriteByte(byte(isByte))

	//Nonce
	b.WriteByte(byte(len(c.Nonce)))
	b.Write(c.Nonce)

	//Cipherbytes
	length := uint32(len(c.Cipherbytes)) // "4GiB ought to be enough for anybody"
	binary.Write(&b, binary.LittleEndian, length)
	b.Write(c.Cipherbytes)

	return b.Bytes(), nil
}

func (c *courier) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)

	//IsCompressed
	isByte, _ := b.ReadByte()
	c.IsCompressed = isByte == 1

	//Nonce
	nlen, _ := b.ReadByte()
	c.Nonce = make([]byte, nlen)
	b.Read(c.Nonce)

	//Cipherbytes
	var clen uint32 // "4GiB ought to be enough for anybody"
	binary.Read(b, binary.LittleEndian, &clen)
	c.Cipherbytes = make([]byte, clen)
	b.Read(c.Cipherbytes)

	return nil
}

type Box struct {
	From string
	Data []byte
}
