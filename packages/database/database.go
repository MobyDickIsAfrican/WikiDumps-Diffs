package database

import (
	"bug/m/submodules/schema"
	"log"
	"os"
	"strconv"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	mu     sync.Mutex
	Table  []schema.DatabaseTable
	Source *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Connect() *Database {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_PATH")), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	d.Source = db
	return d
}

func (d *Database) CreateTable() *Database {
	d.Source.AutoMigrate(&schema.DatabaseTable{})
	return d
}

func (d *Database) Insert(row *schema.DatabaseTable) *Database {
	d.mu.Lock()
	bsc, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		log.Fatal(err)
	}
	if len(d.Table) < bsc {
		d.Table = append(d.Table, *row)
	} else {
		log.Print("Inserting batch")
		d.Source.CreateInBatches(d.Table, bsc)
		d.Table = make([]schema.DatabaseTable, 0)
	}
	d.mu.Unlock()
	return d
}
