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
// type linkInSQLItem struct {
// 	UUID        int    `json:"uuid"`
// 	ShortURL    string `json:"short_url"`
// 	OriginalURL string `json:"original_url"`
// }

// New -.
func NewLinksRepoPGSQL(db *pgx.Conn) *linkRepoInPGSQL {
	return &linkRepoInPGSQL{db}
}

// Insert -.
func (r *linkRepoInPGSQL) Insert(originalURL string) (string, error) {
	shortURL := util.RandomString(cfg.ShortCodeLength)
	_, err := r.db.Exec(context.Background(), `insert into links(original_link, short_link) values ($1, $2)`, originalURL, shortURL)
	if err != nil {
		return shortURL, err
	}
	return shortURL, nil
}

// GetByCode -.
func (r *linkRepoInPGSQL) GetByCode(shortURL string) (string, error) {
	// TODO
	var originalURL string
	err := r.db.QueryRow(context.Background(), "select original_link from links where short_link=$1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", errors.New("link not found")
	}

	return originalURL, nil
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
