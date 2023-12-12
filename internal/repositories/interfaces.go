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

			createTableSQL := `
				DROP TABLE IF EXISTS "public"."links";
				CREATE TABLE "public"."links" (
				"id" int4 NOT NULL GENERATED ALWAYS AS IDENTITY (
				INCREMENT 1
				MINVALUE  1
				MAXVALUE 2147483647
				START 1
				CACHE 1
				),
				"original_link" text COLLATE "pg_catalog"."default" NOT NULL,
				"short_link" text COLLATE "pg_catalog"."default" NOT NULL
				);
				ALTER TABLE "public"."links" ADD CONSTRAINT "videos_pkey" PRIMARY KEY ("id");
				`
			_, dbError := db.Exec(context.Background(), createTableSQL)
			if dbError != nil {
				logrus.Error(dbError)
			} else {
				logrus.Info("Tables created")
			}

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
