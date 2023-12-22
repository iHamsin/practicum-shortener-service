package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/iHamsin/practicum-shortener-service/internal/util"
	"github.com/sirupsen/logrus"
)

var key = "1234567890123456"

func WithCookieCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		UUIDCookie, _ := r.Cookie("UUID")
		UUIDSignCookie, _ := r.Cookie("UUIDSign")
		var UUID string

		if UUIDCookie == nil && UUIDSignCookie == nil {

			UUID, _ = util.GenerateRandomString(10)

			logrus.Debug("New UUID: ", UUID)

			h := hmac.New(sha256.New, []byte(key))
			h.Write([]byte(UUID))
			cryptedNewUUID := h.Sum(nil)

			http.SetCookie(rw, &http.Cookie{
				Name:  "UUID",
				Value: UUID,
			})

			http.SetCookie(rw, &http.Cookie{
				Name:  "UUIDSign",
				Value: hex.EncodeToString(cryptedNewUUID),
			})
			logrus.Debug("New Cookie ", UUID)
		} else if UUIDSignCookie == nil {
			logrus.Debug("No UUID sign ")
			http.Error(rw, "No UUID sign", http.StatusUnauthorized)
			return
		} else {

			logrus.Debug("Check Cookie ", UUIDCookie.Value)

			h := hmac.New(sha256.New, []byte(key))
			h.Write([]byte(UUIDCookie.Value))
			referenceSign := h.Sum(nil)
			userSign, _ := hex.DecodeString(UUIDSignCookie.Value)
			if string(referenceSign) != string(userSign) {
				logrus.Debug("Broken UUID sign")
				http.Error(rw, "Broken UUID sign", http.StatusUnauthorized)
			}
			UUID = UUIDCookie.Value
		}
		ctx := context.WithValue(r.Context(), "UUID", UUID)
		next.ServeHTTP(rw, r.WithContext(ctx))

	})
}
