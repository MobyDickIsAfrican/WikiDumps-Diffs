package hash

import (
	"crypto/md5"
	"fmt"
	"sync"
)

func Hash(text string) string {
	file := []byte(text)

	hash := md5.New()

	// Create a wait group to synchronize the goroutines.
	var wg sync.WaitGroup

	chunkCh := make(chan []byte)
	length := len(file) / 1024
	// Start a goroutine for each chunk of the file.
	for i := 0; i < length; i += 1024 {
		// Get the next chunk of the file.
		chunk := file[i : i+1024]

		go func(chunk []byte) {
			hash.Write(chunk)

			wg.Done()
		}(chunk)
	}

	close(chunkCh)

	wg.Wait()

	hashSum := hash.Sum(nil)
	sum := fmt.Sprintf("%x", hashSum)
	return sum
}
