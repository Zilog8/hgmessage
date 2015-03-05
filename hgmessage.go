package hgmessage

type courier struct {
	Cipherbytes, Nonce []byte
	IsCompressed       bool
	From               string
}

type Box struct {
	From string
	Data []byte
}
