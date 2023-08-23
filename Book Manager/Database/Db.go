package Database

import (
	"Book_Manager/Config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	cfg Config.Config
	db  *gorm.DB
}

func CreateAndConnectToDb(cfg Config.Config) (*GormDB, error) {
	c := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	// Create a new connection
	db, err := gorm.Open(postgres.Open(c), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &GormDB{
		cfg: cfg,
		db:  db,
	}, nil
}
func (gdb *GormDB) CreateModel() error {
	err := gdb.db.AutoMigrate(UserDb{})
	if err != nil {
		return err
	}
	err = gdb.db.AutoMigrate(Book{})
	if err != nil {
		return err
	}
	err = gdb.db.AutoMigrate(Author{})
	if err != nil {
		return err
	}
	return nil

}
