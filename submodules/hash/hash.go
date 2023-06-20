package hash

import (
	"crypto/md5"
	"fmt"
)

func Hash(text string) string {
	file := []byte(text)

	hash := md5.New()

	chunkCh := make(chan []byte)
	length := len(file) / 1024
	for i := 0; i < length; i += 1024 {
		var chunk []byte
		if i+1024 > length {
			chunk = file[i:]
		} else {
			chunk = file[i : i+1024]
		}

		hash.Write(chunk)
	}

	defer close(chunkCh)
	hashSum := hash.Sum(nil)
	return fmt.Sprintf("%x", hashSum)
}
