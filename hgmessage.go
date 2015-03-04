package hgmessage

type Courier struct {
	Cipherbytes, Nonce []byte
	Compressed         bool
}
