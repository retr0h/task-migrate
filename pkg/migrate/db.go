package migrate

import (
	"context"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Version type represents DB model.
type Version struct {
	gorm.Model
	ID   int64  `gorm:"primaryKey;autoIncrement;not null"`
	Name string `gorm:"type:varchar(255);not null; unique"`
}

// Repository contains the DB controller of the Version table.
type Repository struct {
	*gorm.DB
	ctx context.Context
}

// NewRepository creates an instance of Repository.
func NewRepository(
	db *gorm.DB,
) *Repository {
	ctx := context.Background()

	return &Repository{
		DB:  db,
		ctx: ctx,
	}
}

// OpenDB initializes SQLite database connection.
func OpenDB() (*gorm.DB, error) {
	dbPath := "database.db"
	sqliteDialector := sqlite.Open(dbPath)
	db, err := gorm.Open(sqliteDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.AutoMigrate(Version{})
	if err != nil {
		return nil, fmt.Errorf("failed to create db schema: %w", err)
	}

	return db, nil
}

// Insert add a new `Migration` record into the database,
func (r *Repository) Insert(name string) error {
	migration := &Version{
		Name: name,
	}

	err := r.WithContext(r.ctx).Create(migration).Error
	if err != nil {
		return fmt.Errorf("failed to insert migration record: %w", err)
	}

	return nil
}

// Delete remove the given `Migration.Name` from the database,
func (r *Repository) Delete(name string) error {
	var migrations *Version

	err := r.WithContext(r.ctx).
		Where("name = ?", name).
		Delete(&migrations).
		Error
	if err != nil {
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	return nil
}

// FindByName find the given `Migration.Name` from the database,
func (r *Repository) FindByName(name string) (*Version, error) {
	var migrations *Version

	err := r.WithContext(r.ctx).
		Where("name = ?", name).
		First(&migrations).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to find migration record: %w", err)
	}

	return migrations, nil
}

// Find return all records from the database,
func (r *Repository) Find() ([]*Version, error) {
	var migrations []*Version

	err := r.WithContext(r.ctx).Find(&migrations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find migration record: %w", err)
	}

	return migrations, nil
}

// First find the first `Migration` record from the database,
func (r *Repository) First() (*Version, error) {
	var migrations *Version

	err := r.WithContext(r.ctx).First(&migrations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find the first record: %w", err)
	}

	return migrations, nil
}

// Last find the last `Migration` record from the database,
func (r *Repository) Last() (*Version, error) {
	var migrations *Version

	err := r.WithContext(r.ctx).Last(&migrations).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find the last record: %w", err)
	}

	return migrations, nil
}
