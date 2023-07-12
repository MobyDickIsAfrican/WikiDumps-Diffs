package main

import (
	"bufio"
	database "bug/m/packages/database"
	parser "bug/m/packages/parser"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"path/filepath"

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
	//useFolder(os.Getenv("FOLDER_PATH"), file, remaining, dataCh)

	wg.Wait()
}

func useFolder(folder string, file *os.File, remaining chan int64, dataCh chan [][]byte) {
	folderPath := folder

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if !info.IsDir() {
			fmt.Println(path)
			file, err := os.Open(path)
			if err != nil {
				log.Fatal("Error opening file:", err)
			}

			MoveJSONToDatabase(file, remaining, dataCh)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through folder: %v\n", err)
	}

	defer close(remaining)
	defer close(dataCh)
	log.Print("Closing remaining and dataCh")
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

type Stack struct {
	items []interface{}
}

func (s *Stack) Push(item interface{}) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
	if len(s.items) == 0 {
		return nil
	}

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return item
}

func (s *Stack) Peek() interface{} {
	if len(s.items) == 0 {
		return nil
	}

	return s.items[len(s.items)-1]
}

func MoveXMLToDatabase(file *os.File, remaining chan int64, dataCh chan [][]byte) {
	fsInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := float64(fsInfo.Size())
	initSize := fileSize

	var reachedEnd bool = false
	XLMData := &XMLMarshal{ArticleBody: &ArticleBody{}, Version: &Version{}}
	decoder := xml.NewDecoder(file)

	stk := &Stack{}
	for {
		// Read the next token.
		t, err := decoder.Token()
		if reachedEnd {
			reachedEnd = false

			if XLMData.Name == "" && XLMData.Identifier == 0 && XLMData.URL == "" && XLMData.DateModified == "" && XLMData.ArticleBody.Html == "" && XLMData.Version.Identifier == 0 {
				continue
			}
			jsonData, err := json.Marshal(XLMData)
			if err != nil {
				log.Println(err)
			}

			var chk [][]byte
			chk = append(chk, jsonData)
			dataCh <- chk

			remaining <- int64((fileSize / initSize) * 100)
		}

		fileSize -= float64(5)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}

		// set parent.
		if se, ok := t.(xml.StartElement); ok {
			name := se.Name.Local
			stk.Push(name)
		}

		// pop parent
		if _, ok := t.(xml.EndElement); ok {
			stk.Pop()
		}
		if se, ok := t.(xml.CharData); ok {
			currentParent := stk.Pop().(string)

			// skip ns

			if currentParent == "ns" {
				//log.Print("Skipping ns: ", string(se))
				nsp, err := strconv.Atoi(string(se))

				if err != nil {
					log.Fatal(err)
				}

				if nsp != 0 {
					XLMData = &XMLMarshal{ArticleBody: &ArticleBody{}, Version: &Version{}}
					reachedEnd = true
					continue
				}
			}

			if currentParent == "title" {
				XLMData.Name = string(se)
			}

			if currentParent == "id" {
				oldParent := stk.Peek().(string)
				if oldParent == "page" {
					idf, err := strconv.Atoi(string(se))
					if err != nil {
						log.Fatal(err)
					}
					XLMData.Identifier = float64(idf)
				} else if oldParent == "revision" {
					idf, err := strconv.Atoi(string(se))
					if err != nil {
						log.Fatal(err)
					}
					XLMData.Version.Identifier = float64(idf)
				}
			}

			if currentParent == "url" {
				XLMData.URL = string(se)
			}

			if currentParent == "timestamp" {
				XLMData.DateModified = string(se)
			}

			if currentParent == "text" {
				//log.Print("Text: ", string(se))
				XLMData.ArticleBody.Html = "dummyhash"
				reachedEnd = true
			}

			stk.Push(currentParent)
		}
		continue
	}

	defer close(remaining)
	defer close(dataCh)
	log.Print("Closing remaining and dataCh")
}
