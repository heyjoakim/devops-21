package services

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/prometheus"
)

// DBContext defines the application
type DBContext struct {
	db *gorm.DB
}

var (
	dsn         string
	environment string
	dbContext   DBContext
	lock        = &sync.Mutex{}
)

func configureEnv() {
	envFilePath := helpers.GetFullPath("../.env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Println("Error loading .env file - using system variables.")
	}

	environment = os.Getenv("ENVIRONMENT")
	dsn = os.Getenv("DB_CONNECTION")
}

func (d *DBContext) initialize() {
	configureEnv()
	db, err := d.connectDB()
	if err != nil {
		log.Panic(err)
	}

	db.Use(prometheus.New(prometheus.Config{
		DBName:          "Pushgateway",                // use `DBName` as metrics label
		RefreshInterval: 2,                            // Refresh metrics interval (default 15 seconds)
		PushAddr:        "http://142.93.103.26:9091/", // push metrics if `PushAddr` configured
		StartServer:     true,                         // start http server to expose metrics
		HTTPServerPort:  8080,                         // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
		// MetricsCollector: []prometheus.MetricsCollector{
		// 	&prometheus.MySQL{
		// 		VariableNames: []string{"Threads_running"},
		// 	},
		// }, // user defined metrics
	}))

	// db.Callback().Query().Before("gorm:query").Register("my_plugin:before_query", beforeQuery)
	// db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", afterQuery)
	d.db = db
}

var start time.Time

func beforeQuery(db *gorm.DB) {
	start = time.Now()
}

func afterQuery(db *gorm.DB) {
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

func (d *DBContext) connectDB() (*gorm.DB, error) {
	fmt.Println(environment)
	switch environment {
	case "develop":
		fmt.Println("Using local SQLite db")
		return gorm.Open(sqlite.Open("./tmp/minitwit.db"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "production":
		fmt.Println("Using remote postgres db")
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
	case "testing":
		fmt.Println("Using in memory SQLite db")

		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	}
	log.Panic("Environment is not specified in the .env file")
	return nil, nil
}

// initDB creates the database tables.
func (d *DBContext) initDB() {
	err := d.db.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{}, &models.Config{})

	if err != nil {
		log.Println("Migration error:", err)
	}
}

// GetDBInstance returns DBContext with specific environment db
func GetDBInstance() DBContext {
	if dbContext.db == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbContext.db == nil {
			log.Println("Creating Single Instance Now")
			dbContext.initialize()
			dbContext.initDB() // AutoMigrate
		}
	}
	return dbContext
}
