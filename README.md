README
========

hgmessage (or mercury message) is a small and simple wrapper for Go's crypto
and lzma packages, which allows for easy transmission of data that is
compressed (lzma), encrypted (aes) and authenticated (gcm).

Blocking Send. Makes a new connection every time: 

import	"github.com/Zilog8/hgmessage"

err := hgmessage.Send(data, compressionlevel, encryptionkey, recipient)

arguments        | type    | description
---------------- | ------- | ----------------------------------
data             | Message |  interface Message, which implements encoding.BinaryMarshaler
compressionlevel | int     |  LZMA compression level from 1-9. Anything else means no compression.
encryptionkey    | []byte  |  A 128-, 192-, or 256-bit key to encrypt with.
recipient        | string  |  Where to send the data, e.g. "127.0.0.1:4040".

returned         | type    | description
---------------- | ------- | ----------------------------------
err              | error   |  Error if any, else nil.

Buffered Channel Send. Uses a single connection:

import	"github.com/Zilog8/hgmessage"

sendChannel, err := hgmessage.SendChannel(compressionlevel, encryptionkey, recipient)

arguments        | type    | description
---------------- | ------- | ----------------------------------
compressionlevel | int     |  LZMA compression level from 1-9. Anything else means no compression.
encryptionkey    | []byte  |  A 128-, 192-, or 256-bit key to encrypt with.
recipient        | string  |  Where to send the data, e.g. "127.0.0.1:4040".

returned         | type          | description
---------------- | -------------- | ----------------------------------
sendChannel      | chan<- Message |  Accepts Message for transmission.
err              | error         |  Error if any, else nil.

Receiving: 

import	"github.com/Zilog8/hgmessage"

letterChannel, err := hgmessage.ReceiveChannel(encryptionkey, port, senders, mum)


arguments        | type    | description
---------------- | ------------------ | ----------------------------------
encryptionkey    | []byte  |  The key used to encrypt the data.
port             | string  |  Port to receive at, e.g. ":4040".
senders          | string  |  Permited senders; Matches as a prefix. Example: "127.0." matches "127.0.0.1:50437"
mum				 | MessageUnmarshaler |  func([]byte) (Message, error). Unmarshals []byte back into a Message

returned         | type        | description
---------------- | -------------- | ----------------------------------
letterChannel    | <-chan Letter  |  Pumps out Letter{Data: Message, From: string}; The data and who it's from
err              | error       |  Error if any, else nil.

hgmessage
(C) 2014, Zilog8 <zeuscoding@gmail.com>

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.
3. The name of the author may not be used to endorse or promote products
   derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR IMPLIED
WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO
EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS;
OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
