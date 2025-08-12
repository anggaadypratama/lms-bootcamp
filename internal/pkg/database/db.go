package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"fmt"
)

type DatabaseConnection struct {
	DB *gorm.DB
}

func NewDatabaseConnection(dbUrl string) (*DatabaseConnection, error) {
	conn, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	postgres, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	postgres.SetMaxIdleConns(10)
	postgres.SetMaxOpenConns(100)

	return &DatabaseConnection{DB: conn}, nil
}

func (conn *DatabaseConnection) Migrate(models ...interface{}) error {
    if err := conn.DB.AutoMigrate(models...); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }
    log.Println("Migration completed successfully")
    return nil
}


func (dc *DatabaseConnection) Seed() error {
    log.Println("Seeding database...")
    return nil
}

func (dc *DatabaseConnection) Close() error {
    postgres, err := dc.DB.DB()
    if err != nil {
        return err
    }
    return postgres.Close()
}
