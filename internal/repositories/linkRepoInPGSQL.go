package repositories

import (
	"context"
	"errors"

	"github.com/iHamsin/practicum-shortener-service/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// LinkRepoRAM -.
type linkRepoInPGSQL struct {
	db *pgx.Conn
}

// Link -.
type linkInSQLItem struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// New -.
func NewLinksRepoPGSQL(db *pgx.Conn) *linkRepoInPGSQL {
	return &linkRepoInPGSQL{db}
}

// Insert -.
func (r *linkRepoInPGSQL) Insert(originalURL string) (string, error) {
	shortURL := util.RandomString(cfg.ShortCodeLength)

	return shortURL, nil
}

// GetByCode -.
func (r *linkRepoInPGSQL) GetByCode(shortURL string) (string, error) {
	return "", errors.New("link not found")
}

// Close -.
func (r *linkRepoInPGSQL) Close() {
	r.db.Close(context.Background())
}

// Check -.
func (r *linkRepoInPGSQL) Check() error {
	dbError := r.db.Ping(context.Background())
	if dbError != nil {
		return dbError
	}
	logrus.Info("Ping from DB ok")
	return nil
}
