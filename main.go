package main

import (
	"bufio"
	database "bug/m/packages/database"
	parser "bug/m/packages/parser"
	"encoding/json"
	"encoding/xml"
	"io"
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
				prs := parser.ParseXML(data[0])
				dtb.Insert(prs.GetContent())
			}
		}(&dataCh, wg, dtb)
	}

	MoveXMLToDatabase(file, remaining, dataCh)

	wg.Wait()
}

func MoveJSONToDatabase(file *os.File, remaining chan int64, dataCh chan [][]byte) {
	fsInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := float64(fsInfo.Size())

	initSize := fileSize

	scanner := bufio.NewScanner(file)
	const maxBufferSize = 1024 * 1024
	scanner.Buffer(make([]byte, maxBufferSize), maxBufferSize)
	for scanner.Scan() {
		log.Print(time.Now())
		var chk [][]byte

		line := scanner.Text()
		log.Print("Line: ", line)
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
}

type XMLMarshal struct {
	Name         string       `json:"name"`
	Identifier   float64      `json:"identifier"`
	Version      *Version     `json:"version"`
	URL          string       `json:"url"`
	DateModified string       `json:"date_modified"`
	ArticleBody  *ArticleBody `json:"article_body"`
}
type ArticleBody struct {
	Html string `json:"html"`
}

type Version struct {
	Identifier float64 `json:"identifier"`
}

func MoveXMLToDatabase(file *os.File, remaining chan int64, dataCh chan [][]byte) {
	fsInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := float64(fsInfo.Size())
	initSize := fileSize

	var reachedEnd bool = false
	XLMData := XMLMarshal{ArticleBody: &ArticleBody{}, Version: &Version{}}
	decoder := xml.NewDecoder(file)
	// Iterate over the XML tags.
	for {
		// Read the next token.
		t, err := decoder.Token()
		if reachedEnd {
			reachedEnd = false
			jsonData, err := json.Marshal(XLMData)
			log.Printf("JSON: %s", jsonData)
			if err != nil {
				log.Println(err)
			}

			var chk [][]byte
			chk = append(chk, jsonData)
			dataCh <- chk

			remaining <- int64((fileSize / initSize) * 100)
			fileSize -= float64(len(jsonData))
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}

		// Check if the token is a StartElement.
		if se, ok := t.(xml.StartElement); ok {
			// Get the name of the tag.
			name := se.Name.Local
			value := se.Attr[0].Value

			if name == "title" {
				XLMData.Name = value
			}
			if name == "id" {
				XLMData.Identifier = 1
				XLMData.Version.Identifier = 1
			}
			/*
				if name == "version_identifier" {
					XLMData.VersionIdentifier = 1
				}*/

			if name == "url" {
				XLMData.URL = name
			}

			if name == "timestamp" {
				XLMData.DateModified = name
			}

			if name == "text" {
				XLMData.ArticleBody.Html = name
				reachedEnd = true
			}

			continue
		}
		continue
	}

	defer close(remaining)
	defer close(dataCh)
	log.Print("Closing remaining and dataCh")
}
