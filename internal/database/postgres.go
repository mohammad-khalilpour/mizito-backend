package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mizito/internal/env"
	"mizito/pkg/models"
)

type DatabaseHandler struct {
	DB    *gorm.DB
	Cfg   *env.Config
	Redis RedisHandler // Optional, depending on your needs
}

func createPostgresDSN(cfg *env.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDatabase,
	)
}

// NewDatabaseHandler initializes the PostgreSQL database handler.
func NewDatabaseHandler(cfg *env.Config) *DatabaseHandler {
	var dbHandler DatabaseHandler
	// Initialize the PostgreSQL connection
	db, err := gorm.Open(postgres.Open(createPostgresDSN(cfg)), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to PostgreSQL: %s", err))
	}

	// Set up the handler
	dbHandler = DatabaseHandler{
		DB:  db,
		Cfg: cfg,
	}

	// Ping the database to ensure the connection is live
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve database connection: %s", err))
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("failed to ping database: %s", err))
	}

	// Run migrations
	if err := dbHandler.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %s", err))
	}

	return &dbHandler
}

// Migrate runs database migrations for all models.
func (d *DatabaseHandler) Migrate() error {
	if err := d.DB.AutoMigrate(
		&models.Team{},
		&models.Project{},
		&models.Task{},
		&models.Subtask{},
		&models.User{},
		&models.TeamMember{},
		&models.Message{},
		&models.Report{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
