package main

import (
	"bufio"
	parser "bug/m/packages/parser"
	"log"
	"os"
	"strings"
	"sync"
)

type Data struct {
	// Define the structure of your JSON data
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
	// ...
}

func worker(ch <-chan [][]byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for dataChunk := range ch {
		parser.ParseJSON(dataChunk[0])
	}
}

func main() {
	// Open the JSON file
	file, err := os.Open("./packages/parser/testdata/test_data.ndjson")
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// Create a channel for sending JSON data chunks to workers
	dataCh := make(chan [][]byte, 6)

	// Start worker goroutines
	var wg sync.WaitGroup
	numWorkers := 5
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker(dataCh, &wg)
	}

	// Read the JSON file in chunks and send them to workers
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//dataChunk := scanner.Bytes()
		//log.Print(scanner.Text())
		//var chk [][]byte
		//chk = append(chk, dataChunk)
		// Send the data chunk to the worker goroutines for processing
		//dataCh <- chk
		var chk [][]byte

		line := scanner.Text()

		// Escape the URL characters in the line.
		line = strings.ReplaceAll(line, "&", "&amp;")
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		//log.Print(line)
		chk = append(chk, []byte(line))
		dataCh <- chk
	}

	// Close the channel to indicate no more data will be sent
	close(dataCh)

	// Wait for all workers to complete
	wg.Wait()
}
