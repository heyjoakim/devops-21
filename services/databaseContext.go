package services

import (
	"fmt"
	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"sync"
)

// DbContext defines the application
type DbContext struct {
	db *gorm.DB
}

var dsn string
var dbContext DbContext
var lock = &sync.Mutex{}

func configureEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file - using system variables.")
	}
	dsn = os.Getenv("DB_CONNECTION")
}

func (d *DbContext) init() {
	configureEnv()
	db, err := d.connectDb()
	if err != nil {
		log.Panic(err)
	}
	d.db = db
}

func (d *DbContext) connectDb() (*gorm.DB, error) {
	return gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
}

// initDb creates the database tables.
func (d *DbContext) initDb() {
	err := d.db.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{})
	if err != nil {
		fmt.Println("Migration error:", err)
	}
}

func GetDbInstance() DbContext {
	if dbContext.db != nil {
		fmt.Println("Single Instance already created-2")
	} else {
		lock.Lock()
		defer lock.Unlock()
		if dbContext.db != nil {
			fmt.Println("Single Instance already created-1")
		} else {
			fmt.Println("Creating Single Instance Now")
			dbContext.init()
			dbContext.initDb() // AutoMigrate
		}
	}
	return dbContext

}
