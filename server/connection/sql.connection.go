package connection

import (
	"backend/server/env"
	"backend/server/initialize"
	"backend/server/structure"
	"backend/server/util"
	"fmt"
	"log"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQL struct{}

var opts = &gorm.Config{
	// options...
}

var isConnected = false

func (ref SQL) Connect() (*gorm.DB, error) {
	dsn := structure.GormDSN{
		Host: env.GetDbHost(),
		Port: env.GetDbPort(),
		User: env.GetDbUser(),
		Pass: env.GetDbPass(),
		Name: env.GetDbName(),
	}
	db_type := env.GetDbType()
	var db *gorm.DB
	var err error
	if db_type == "mysql" {
		db, err = ref.MySQL(dsn)
	} else if db_type == "postgres" {
		db, err = ref.PostgreSQL(dsn)
	} else if db_type == "mssql" {
		db, err = ref.MSSQL(dsn)
	} else {
		db, err = ref.MySQL(dsn)
	}
	if err != nil {
		log.Fatalf("❌ Database error: %s", err)
	}
	if !isConnected {
		log.Println("✅ Database Connected")

		if env.GetDbMigration() {
			if err := db.AutoMigrate(initialize.AutoMigrationTables()...); err != nil {
				panic("failed to auto migrate database: " + err.Error())
			}
			log.Println("✅ Auto Migration")
		} else {
			log.Println("👌 Auto Migration Skipped")
		}

		isConnected = true
	}
	return db, err
}

// ---------------------------------------------------------------------------------------------------------------

func (ref SQL) SQLite(dsn structure.GormDSN) (*gorm.DB, error) {
	File := util.File{}
	pwd := env.GetPwd()
	databaseDir := filepath.Join(pwd, "database")
	databaseName := File.AddExtensionIfNotExist(dsn.Name, "db")
	databasePath := filepath.Join(databaseDir, databaseName)
	err := File.CreateIfNotExist(databasePath)
	if err != nil {
		return nil, fmt.Errorf("error create database file: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(databasePath), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (ref SQL) SQLiteInMemory(dsn structure.GormDSN) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ref SQL) MySQL(dsn structure.GormDSN) (*gorm.DB, error) {
	dsn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dsn.User, dsn.Pass, dsn.Host, dsn.Port, dsn.Name)
	db, err := gorm.Open(mysql.Open(dsn_str), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ref SQL) PostgreSQL(dsn structure.GormDSN) (*gorm.DB, error) {
	dsn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dsn.Host, dsn.Port, dsn.User, dsn.Pass, dsn.Name)
	db, err := gorm.Open(postgres.Open(dsn_str), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (ref SQL) MSSQL(dsn structure.GormDSN) (*gorm.DB, error) {
	dsn_str := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", dsn.Host, dsn.User, dsn.Pass, dsn.Port, dsn.Name)
	db, err := gorm.Open(sqlserver.Open(dsn_str), opts)
	if err != nil {
		return nil, err
	}
	return db, nil
}
