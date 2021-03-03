package services

import (
	"fmt"
	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"sync"
)

// DbContext defines the application
type DbContext struct {
	DB *gorm.DB
}

var (
	dsn         string
	environment string
	dbContext   DbContext
	lock        = &sync.Mutex{}
)

func configureEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("hep")
		log.Println("Error loading .env file - using system variables.")
	}

	environment = os.Getenv("ENVIRONMENT")
	dsn = os.Getenv("DB_CONNECTION")
}

func (d *DbContext) initialize() {
	configureEnv()
	db, err := d.connectDb()
	if err != nil {
		log.Panic(err)
	}
	d.DB = db
}

/*
import "github.com/joho/godotenv"
func TestSendMessage(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error("Error loading .env file")
	}
}

 */

func (d *DbContext) connectDb() (*gorm.DB, error) {
	fmt.Println(environment)
	if environment == "develop" {
		fmt.Println("Using local SQLite DB")
		return gorm.Open(sqlite.Open("./tmp/minitwit.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	} else if environment == "production" {
		fmt.Println("Using remote postgres DB")
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
	} else if environment == "testing" {
		fmt.Println(environment)
		fmt.Println("Using in memory SQLite DB")
		return gorm.Open(sqlite.Open("./tmp/test.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	}
	log.Panic("Environment is not specified in the .env file")
	return nil, nil
}

// initDb creates the database tables.
func (d *DbContext) initDb() {
	err := d.DB.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{})
	if err != nil {
		log.Println("Migration error:", err)
	}
}

func GetDbInstance() DbContext {
	if dbContext.DB == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbContext.DB == nil {
			log.Println("Creating Single Instance Now")
			dbContext.initialize()
			dbContext.initDb() // AutoMigrate
		}
	}
	return dbContext

}
