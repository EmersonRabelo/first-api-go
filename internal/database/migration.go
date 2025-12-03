package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunMigrations executa as migrations do banco de dados
func RunMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter instância sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao criar driver de migrations: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar instância de migrations: %w", err)
	}

	// Executa as migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("erro ao executar migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("erro ao obter versão da migration: %w", err)
	}

	if dirty {
		log.Printf("ATENÇÃO: Banco de dados em estado dirty (versão %d)", version)
	} else if err == nil {
		log.Printf("Migrations executadas com sucesso! Versão atual: %d", version)
	} else {
		log.Println("Nenhuma migration foi executada ainda")
	}

	return nil
}

// RollbackMigration reverte a última migration
func RollbackMigration(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter instância sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao criar driver de migrations: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar instância de migrations: %w", err)
	}

	if err := m.Steps(-1); err != nil {
		return fmt.Errorf("erro ao reverter migration: %w", err)
	}

	version, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("erro ao obter versão da migration: %w", err)
	}

	if err == nil {
		log.Printf("Migration revertida com sucesso! Versão atual: %d", version)
	} else {
		log.Println("Todas as migrations foram revertidas")
	}

	return nil
}

// MigrationStatus retorna o status atual das migrations
func MigrationStatus(db *gorm.DB) (uint, bool, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return 0, false, fmt.Errorf("erro ao obter instância sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("erro ao criar driver de migrations: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("erro ao criar instância de migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, fmt.Errorf("erro ao obter versão da migration: %w", err)
	}

	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}

	return version, dirty, nil
}

// CreateMigrationDatabase cria o banco de dados se não existir
func CreateMigrationDatabase(host, port, user, password, dbName string) error {
	// Conecta ao postgres database padrão para criar o banco
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao PostgreSQL: %w", err)
	}
	defer db.Close()

	// Verifica se o banco existe
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("erro ao verificar se banco existe: %w", err)
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("erro ao criar banco de dados: %w", err)
		}
		log.Printf("Banco de dados '%s' criado com sucesso!", dbName)
	} else {
		log.Printf("Banco de dados '%s' já existe", dbName)
	}

	return nil
}
