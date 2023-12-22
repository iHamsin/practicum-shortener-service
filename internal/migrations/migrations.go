package migrations

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sirupsen/logrus"
)

//go:embed *.sql
var fs embed.FS

func MigrationsUP(databaseDSN string) error {

	var outError error

	logrus.Info("Migrations started")

	d, err := iofs.New(fs, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseDSN)

	if err != nil {
		logrus.Error("failed to get a new migrate instance: ", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logrus.Error("failed to apply migrations to the DB: ", err)
		}
	}

	logrus.Info("Migrations complited")
	return outError
}
