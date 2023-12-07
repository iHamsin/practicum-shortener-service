package repositories

import (
	"os"

	"github.com/iHamsin/practicum-shortener-service/config"
	"github.com/sirupsen/logrus"
)

var cfg *config.Config

type Repository interface {
	// GetAll() (map[string]string, error)
	GetByCode(string) (string, error)
	Insert(string) (string, error)
	Close()
}

func Init(incomeCfg *config.Config) (Repository, error) {
	cfg = incomeCfg
	var repository Repository
	if cfg.HTTP.DBFile == "" {
		logrus.Debug("DB in RAM")
		repository = NewLinksRepoRAM(make(map[string]string))
	} else {
		logrus.Debug("DB in file", cfg.HTTP.DBFile)

		file, fileOpenError := os.OpenFile(cfg.HTTP.DBFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if fileOpenError != nil {
			logrus.Error(fileOpenError)
		}

		var fileRepoError error
		repository, fileRepoError = NewLinksRepoFile(*file)
		if fileRepoError != nil {
			logrus.Error(fileRepoError)
		}
	}
	return repository, nil
}
