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
func (r *linksRepoInRAM) Insert(link string) (string, error) {
	// генерируем ключ и проверяем на наличие такого в хранилище
	var linkKey string
	var i = 0
	for {
		linkKey = util.RandomString(8)
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
	r.storage[linkKey] = link
	return linkKey, nil
}

// GetByCode -.
func (r *linksRepoInRAM) GetByCode(linkCode string) (string, error) {
	// проверка наличия в хранилище
	_, URLfound := r.storage[linkCode]
	if !URLfound {
		return "", errors.New("link not found")
	}
	return r.storage[linkCode], nil
}
