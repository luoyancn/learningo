package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5Crypto(context string) string {
	md5_writer := md5.New()
	io.WriteString(md5_writer, context)
	return hex.EncodeToString(md5_writer.Sum(nil))
}
