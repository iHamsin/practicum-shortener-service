package repositories

import (
	"context"
	"embed"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

var cfg *config.Config

type Repository interface {
	// GetAll() (map[string]string, error)
	GetLinkByCode(string) (string, error)
	InsertLink(string) (string, error)
	Check() error
	Close()
	BatchInsertLink([]string) ([]string, error)
}

//go:embed migrations/*.sql
var fs embed.FS

func Init(incomeCfg *config.Config) (Repository, error) {
	cfg = incomeCfg
	var repository Repository
	var outError error
	if cfg.Repository.DatabaseDSN != "" {
		logrus.Info("DB in postgres", cfg.Repository.DatabaseDSN)

		d, err := iofs.New(fs, "migrations")
		if err != nil {
			return nil, err
		}

		m, err := migrate.NewWithSourceInstance("iofs", d, cfg.Repository.DatabaseDSN)
		if err != nil {
			logrus.Error("failed to get a new migrate instance: ", err)
		}
		if err := m.Up(); err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				logrus.Error("failed to apply migrations to the DB: ", err)
			}
		}

		db, postgresOpenError := pgx.Connect(context.Background(), cfg.Repository.DatabaseDSN)
		if postgresOpenError != nil {
			outError = postgresOpenError
			repository = nil
		} else {
			repository = NewLinksRepoPGSQL(db)
		}
	} else if cfg.HTTP.DBFile != "" {
		logrus.Debug("DB in file", cfg.HTTP.DBFile)

		file, fileOpenError := os.OpenFile(cfg.HTTP.DBFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if fileOpenError != nil {
			logrus.Error(fileOpenError)
		}

		var fileRepoError error
		repository, fileRepoError = NewLinksRepoFile(*file)
		if fileRepoError != nil {
			repository = nil
			outError = fileRepoError
		}
	} else {
		// DB in RAM
		logrus.Debug("DB in RAM")
		repository = NewLinksRepoRAM(make(map[string]string))
		outError = nil
	}
	return repository, outError
}
