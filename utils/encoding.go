package utils

import "golang.org/x/text/encoding/unicode"

func UTF16ToUTF8(byt []byte) ([]byte, error) {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	return decoder.Bytes(byt[:])
}
