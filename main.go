package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/joho/godotenv"
)


type Product struct {
	tableName   struct{}  `pg:"products"`
	ID          int64     `pg:",pk"`
	ProductName string    `pg:",notnull"`
	FarmerID    int64     `pg:",notnull"` 
	CreatedAt   time.Time `pg:"default:now()"`
}

type Stakeholder struct {
	tableName     struct{} `pg:"stakeholders"`
	ID            int64    `pg:",pk"`
	Name          string   `pg:",notnull"`
	Role          string   `pg:",notnull"` 
	WalletAddress string   `pg:",unique"`
}

type SupplyChainEvent struct {
	tableName     struct{}  `pg:"supply_chain_events"`
	ID            int64     `pg:",pk"`
	ProductID     int64     `pg:",notnull"` 
	StakeholderID int64     `pg:",notnull"` 
	Status        string    `pg:",notnull"` 
	Location      string    `pg:",null"`
	Timestamp     time.Time `pg:"default:now()"`
}


func createDBConnection() *pg.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables.")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	
	if dbUser == "" || dbPassword == "" || dbHost == "" {
		log.Fatal("Database environment variables DB_USER, DB_PASSWORD, and DB_HOST must be set.")
	}

	db := pg.Connect(&pg.Options{
		User:     dbUser,
		Password: dbPassword,
		Addr:     fmt.Sprintf("%s:%s", dbHost, dbPort),
		Database: dbName,
	})

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}


func createTables(db *pg.DB) {
	models := []interface{}{
		(*Product)(nil),
		(*Stakeholder)(nil),
		(*SupplyChainEvent)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
			FKConstraints: true, 
		})
		if err != nil {
			log.Fatalf("Error creating table for %T: %v", model, err)
		}
		fmt.Printf("Table for %T created successfully.\n", model)
	}
}


func main() {
	db := createDBConnection()
	defer db.Close()

	createTables(db)
}