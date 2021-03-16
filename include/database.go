package include

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func initDatabase() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	db, dbErr = gorm.Open(sqlite.Open(os.Getenv("DATABASE_PATH")+"python-script-sch.db"), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if dbErr != nil {
		Log.Panic("[DATABASE INIT] ", dbErr)
		panic("ERROR failed to connect database ")
	}

	dbErr = db.AutoMigrate(
		&DBScriptFile{},
		&DBBaseScript{},
		&DBScriptParameter{},
		&DBJob{},
		&DBTaskParameter{},
		&DBTask{},
		&DBReport{},
		&DBRecipient{},
		&DBOutgoingMails{},
		&DBOutgoingMailHistory{},
	)

	if dbErr != nil {
		Log.Panic("[HANDLE DB] ERROR, DB AutoMigrate ", dbErr)
	}

}
