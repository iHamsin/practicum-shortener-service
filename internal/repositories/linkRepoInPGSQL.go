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
	db                   *pgx.Conn
	linksToDeleteChannle chan []string
}

var ErrDublicateOriginalLink = errors.New("dublicate orgiginal link")

// New -.
func NewLinksRepoPGSQL(db *pgx.Conn) *linkRepoInPGSQL {
	newPGSQLRepo := &linkRepoInPGSQL{db, make(chan []string)}
	go func(db *pgx.Conn) {
		for {
			linksToDelete := <-newPGSQLRepo.linksToDeleteChannle
			batch := &pgx.Batch{}
			for _, link := range linksToDelete {
				batch.Queue("UPDATE links SET deleted_flag = TRUE WHERE short_link = $1", link)
			}
			batchResult := db.SendBatch(context.Background(), batch)
			batchResult.Close()
		}
	}(newPGSQLRepo.db)
	return newPGSQLRepo
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

// BatchDeleteLink -.
func (r *linkRepoInPGSQL) BatchDeleteLink(ctx context.Context, links []string, UUID string) (bool, error) {
	r.linksToDeleteChannle <- links
	return true, nil
}

// GetLinkByCode -.
func (r *linkRepoInPGSQL) GetLinkByCode(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	var deletedFlag bool
	err := r.db.QueryRow(ctx, "select original_link, deleted_flag from links where short_link=$1", shortURL).Scan(&originalURL, &deletedFlag)
	if err != nil {
		logrus.Error(err)
		return "", errors.New("link not found")
	}
	if deletedFlag {
		return "", errors.New("link gone")
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
		err := rows.Scan(&link.OriginalURL, &link.ShortURL)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
