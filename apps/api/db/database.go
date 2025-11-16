package db

import (
	"fmt"
	"log"

	"github.com/keeee21/commit-town/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection
func NewDatabase(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate all models in order
	err := db.AutoMigrate(
		&models.User{},
		&models.UserRepository{},
		&models.RepoDailyCommitLog{},
		&models.UserDailyCommitLog{},
		&models.UserStreak{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Add unique constraints
	if err := addUniqueConstraints(db); err != nil {
		return fmt.Errorf("failed to add unique constraints: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// addUniqueConstraints adds unique constraints that are not directly supported by GORM tags
func addUniqueConstraints(db *gorm.DB) error {
	// UserRepository: unique constraint on (UserID, RepoOwner, RepoName)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_repositories_user_repo
		ON user_repositories(user_id, repo_owner, repo_name)
	`).Error; err != nil {
		return err
	}

	// RepoDailyCommitLog: unique constraint on (UserRepoID, CommitDate)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_repo_daily_commit_logs_repo_date
		ON repo_daily_commit_logs(user_repo_id, commit_date)
	`).Error; err != nil {
		return err
	}

	// UserDailyCommitLog: unique constraint on (UserID, Date)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_daily_commit_logs_user_date
		ON user_daily_commit_logs(user_id, date)
	`).Error; err != nil {
		return err
	}

	return nil
}
