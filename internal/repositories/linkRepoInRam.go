package repositories

import (
	"errors"

	"github.com/iHamsin/practicum-shortener-service/internal/util"
)

// LinkRepoRAM -.
type linksRepoInRAM struct {
	storage map[string]string
}

// New -.
func NewLinksRepoRAM(storage map[string]string) *linksRepoInRAM {
	return &linksRepoInRAM{storage}
}

// Insert -.
func (r *linksRepoInRAM) Insert(originalURL string) (string, error) {
	// генерируем ключ и проверяем на наличие такого в хранилище
	var linkKey string
	var i = 0
	for {
		linkKey = util.RandomString(cfg.ShortCodeLength)
		_, keyUsed := r.storage[linkKey]
		if !keyUsed {
			break
		}
		i++
		// если за 20 попыток не сгенерировали униклаьный ключ - сдаемся
		if i > 20 {
			return "", errors.New("cant generate uniq link short code")
		}
	}
	r.storage[linkKey] = originalURL
	return linkKey, nil
}

// BatchInsert -.
func (r *linksRepoInRAM) BatchInsert(links []string) ([]string, error) {
	result := make([]string, len(links))

	for i, link := range links {
		result[i] = util.RandomString(cfg.ShortCodeLength)
		r.storage[result[i]] = link
	}

	return result, nil
}

// GetByCode -.
func (r *linksRepoInRAM) GetByCode(shortURL string) (string, error) {
	// проверка наличия в хранилище
	_, URLfound := r.storage[shortURL]
	if !URLfound {
		return "", errors.New("link not found")
	}
	return r.storage[shortURL], nil
}

// Close -.
func (r *linksRepoInRAM) Close() {

}

// Check -.
func (r *linksRepoInRAM) Check() error {
	return nil
}
