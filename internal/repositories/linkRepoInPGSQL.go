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
	result, err := r.db.Exec(context.Background(), `insert into links(original_link, short_link) values ($1, $2) ON CONFLICT (original_link) DO NOTHING`, originalURL, shortURL)
	if err != nil {
		return "", err
	}
	if result.RowsAffected() == 0 {
		shortURL, err = r.GetByOriginalLink(originalURL)
		if err != nil {
			return "", err
		}
	}
	return shortURL, nil
}

// BatchInsert -.
func (r *linkRepoInPGSQL) BatchInsert(links []string) ([]string, error) {
	result := make([]string, len(links))

	tx, err := r.db.Begin(context.TODO())
	if err != nil {
		return nil, err
	}

	for i, link := range links {
		result[i] = util.RandomString(cfg.ShortCodeLength)
		// insert
		_, err = tx.Exec(context.TODO(), `insert into links(original_link, short_link) values ($1, $2)`, link, result[i])
		if err != nil {
			// если ошибка, то откатываем изменения
			_ = tx.Rollback(context.TODO())
			return nil, err
		}
	}
	_ = tx.Commit(context.TODO())
	return result, nil
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

// GetByCode -.
func (r *linkRepoInPGSQL) GetByOriginalLink(originalLink string) (string, error) {
	// TODO
	var shortLink string
	err := r.db.QueryRow(context.Background(), "select short_link from links where original_link=$1", originalLink).Scan(&shortLink)
	if err != nil {
		return "", errors.New("link not found")
	}

	return shortLink, nil
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
