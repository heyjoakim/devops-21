package services

import (
	"os"
	"sync"

	"github.com/heyjoakim/devops-21/helpers"
	"github.com/heyjoakim/devops-21/models"
	"github.com/joho/godotenv"
	gorm_logrus "github.com/onrik/gorm-logrus"
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
		LogError(models.Log{
			Message: err.Error(),
		})
	}

	environment = os.Getenv("ENVIRONMENT")
	dsn = os.Getenv("DB_CONNECTION")
}

func (d *DBContext) initialize() {
	configureEnv()
	db, err := d.connectDB()
	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
		})
	}

	var (
		logInterval = uint32(2)    //nolint
		HTTPPort    = uint32(8080) //nolint
	)

	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "Pushgateway",                     // use `DBName` as metrics label
		RefreshInterval: logInterval,                       // Refresh metrics interval (default 15 seconds)
		PushAddr:        os.Getenv("PROMETHEUS_PUSH_ADDR"), // push metrics if `PushAddr` configured
		StartServer:     true,                              // start http server to expose metrics
		HTTPServerPort:  HTTPPort,                          // configure http server port, default port 8080
		// (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
	}))

	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
		})
	}

	d.db = db
}

func (d *DBContext) connectDB() (*gorm.DB, error) {
	// debug.PrintStack()

	switch environment {
	case "develop":
		LogInfo("Using local SQLite db")
		return gorm.Open(sqlite.Open("./tmp/minitwit.db"), &gorm.Config{
			Logger: gorm_logrus.New(),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "production":
		LogInfo("Using remote postgres db")
		return gorm.Open(postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			}),
			&gorm.Config{
				Logger: gorm_logrus.New(),
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
			})
	case "testing":
		LogInfo("Using in memory SQLite db")

		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
			Logger: gorm_logrus.New(),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	}
	LogPanic("Environment is not specified in the .env file")
	return nil, nil
}

// initDB creates the database tables.
func (d *DBContext) initDB() {
	err := GetDBInstance().db.AutoMigrate(&models.User{}, &models.Follower{}, &models.Message{}, &models.Config{})
	if err != nil {
		LogFatal(models.Log{
			Message: "Migration error",
			Data:  err,
		})
	}
}

// GetDBInstance returns DBContext with specific environment db
func GetDBInstance() DBContext {
	if dbContext.db == nil {
		lock.Lock()
		defer lock.Unlock()
		if dbContext.db == nil {
			LogInfo("Creating Single Instance Now")
			dbContext.initialize()
			dbContext.initDB() // AutoMigrate
		}
	}
	return dbContext
}
