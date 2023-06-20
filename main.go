package main

import (
	"bufio"
	database "bug/m/packages/database"
	parser "bug/m/packages/parser"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	file, err := os.Open(os.Getenv("FILE_PATH"))
	if err != nil {
		log.Fatal("Error opening file:", err)
	}

	defer file.Close()

	remaining := make(chan int64, 1)
	fsInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := float64(fsInfo.Size())

	initSize := fileSize

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(r *chan int64) {
		defer wg.Done()
		for status := range *r {
			log.Printf("%d%% remaining", status)
		}
	}(&remaining)

	numWorkers := 10
	wg.Add(numWorkers)

	dataCh := make(chan [][]byte)
	dtb := database.NewDatabase().Connect().CreateTable()
	for i := 0; i < numWorkers; i++ {
		go func(r *chan [][]byte, wg *sync.WaitGroup, dtb *database.Database) {
			defer wg.Done()
			for data := range *r {
				prs := parser.ParseJSON(data[0])
				dtb.Insert(prs.GetContent())
			}
		}(&dataCh, wg, dtb)
	}

	scanner := bufio.NewScanner(file)
	const maxBufferSize = 1024 * 1024
	scanner.Buffer(make([]byte, maxBufferSize), maxBufferSize)
	for scanner.Scan() {
		log.Print(time.Now())
		var chk [][]byte

		line := scanner.Text()
		line = strings.ReplaceAll(line, "&", "&amp;")
		line = strings.ReplaceAll(line, "<", "&lt;")
		line = strings.ReplaceAll(line, ">", "&gt;")
		chk = append(chk, []byte(line))

		remaining <- int64((fileSize / initSize) * 100)
		dataCh <- chk
		fileSize -= float64(len(line))
	}

	defer close(remaining)
	defer close(dataCh)
	log.Print("Closing remaining and dataCh")

	wg.Wait()
}
