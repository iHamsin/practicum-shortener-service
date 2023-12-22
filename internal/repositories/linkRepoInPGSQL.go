package repositories

import (
	"context"
	"errors"

	"github.com/iHamsin/practicum-shortener-service/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// linkRepoInPGSQL -.
type linkRepoInPGSQL struct {
	db *pgx.Conn
}

var ErrDublicateOriginalLink = errors.New("dublicate orgiginal link")

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

// InsertLink -.
func (r *linkRepoInPGSQL) InsertLink(ctx context.Context, originalURL string, UUID string) (string, error) {
	shortURL := util.RandomString(cfg.ShortCodeLength)
	result, err := r.db.Exec(ctx, `insert into links(original_link, short_link, uuid) values ($1, $2, $3) ON CONFLICT (original_link) DO NOTHING`, originalURL, shortURL, UUID)
	if err != nil {
		return "", err
	}
	if result.RowsAffected() == 0 {
		shortURL, err = r.GetLinkByOriginalLink(ctx, originalURL)
		if err != nil {
			return "", err
		}
		return shortURL, ErrDublicateOriginalLink
	}
	return shortURL, nil
}

// BatchInsertLink -.
func (r *linkRepoInPGSQL) BatchInsertLink(ctx context.Context, links []string, UUID string) ([]string, error) {
	result := make([]string, len(links))

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	for i, link := range links {
		result[i] = util.RandomString(cfg.ShortCodeLength)
		// insert
		command, err := tx.Exec(ctx, `insert into links(original_link, short_link, uuid) values ($1, $2, $3) ON CONFLICT (original_link) DO NOTHING`, link, result[i], UUID)
		if err != nil {
			// если ошибка, то откатываем изменения
			_ = tx.Rollback(ctx)
			return nil, err
		}
		if command.RowsAffected() == 0 {
			shortURL, err := r.GetLinkByOriginalLink(ctx, link)
			if err != nil {
				_ = tx.Rollback(ctx)
				return nil, err
			}
			result[i] = shortURL
		}
	}
	_ = tx.Commit(context.TODO())
	return result, nil
}

// GetLinkByCode -.
func (r *linkRepoInPGSQL) GetLinkByCode(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := r.db.QueryRow(ctx, "select original_link from links where short_link=$1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", errors.New("link not found")
	}

	return originalURL, nil
}

// GetLinkByOriginalLink -.
func (r *linkRepoInPGSQL) GetLinkByOriginalLink(ctx context.Context, originalLink string) (string, error) {
	var shortLink string
	err := r.db.QueryRow(ctx, "select short_link from links where original_link=$1", originalLink).Scan(&shortLink)
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

// GetLinksByUUID -.
func (r *linkRepoInPGSQL) GetLinksByUUID(ctx context.Context, UUID string) ([]Link, error) {
	var links []Link
	var link Link
	rows, _ := r.db.Query(ctx, "select original_link, short_link from links where uuid=$1", UUID)
	for rows.Next() {
		err := rows.Scan(&link.OriginalLink, &link.ShortLink)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
