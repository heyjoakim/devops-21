package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
    _, b, _, _ = runtime.Caller(0)
    basepath   = filepath.Dir(b)
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
	envFilePath := getFullPath("../.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("hep")
		log.Println("Error loading .env file - using system variables.")
	}

	environment = os.Getenv("ENVIRONMENT")
	dsn = os.Getenv("DB_CONNECTION")
}

func getFullPath(fileName string) string {
		_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env.test)")
		os.Exit(1)
	}
	basepath := filepath.Dir(file)
	return filepath.Join(basepath,fileName)
}

func (d *DbContext) initialize() {
	configureEnv()
	db, err := d.connectDb()
	if err != nil {
		log.Panic(err)
	}
	d.DB = db
}

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
		dbPath := getFullPath("../tmp/minitwit.db")
		fmt.Println(dbPath)
		return gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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
	err := d.DB.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{},&models.Config{})
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
