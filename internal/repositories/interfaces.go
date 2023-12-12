package repositories

import (
	"context"
	"os"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

var cfg *config.Config

type Repository interface {
	// GetAll() (map[string]string, error)
	GetByCode(string) (string, error)
	Insert(string) (string, error)
	Check() error
	Close()
}

func Init(incomeCfg *config.Config) (Repository, error) {
	cfg = incomeCfg
	var repository Repository
	var outError error
	if cfg.HTTP.DBFile == "" && cfg.Repository.DatabaseDSN != "" {
		logrus.Info("DB in postgres", cfg.Repository.DatabaseDSN)

		db, postgresOpenError := pgx.Connect(context.Background(), cfg.Repository.DatabaseDSN)
		if postgresOpenError != nil {
			outError = postgresOpenError
			repository = nil
		} else {
			repository = NewLinksRepoPGSQL(db)
		}
	} else if cfg.HTTP.DBFile != "" && cfg.Repository.DatabaseDSN == "" {
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
