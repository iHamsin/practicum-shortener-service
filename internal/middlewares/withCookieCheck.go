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

type RequestUUIDKey struct{}
type RequestisNewUUIDKey struct{}

func WithCookieCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		UUIDCookie, _ := r.Cookie("UUID")
		UUIDSignCookie, _ := r.Cookie("UUIDSign")
		var UUID string
		var isNewUUID bool

		if UUIDCookie == nil && UUIDSignCookie == nil {

			UUID, _ = util.GenerateRandomString(10)
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
			isNewUUID = true
			logrus.Debug("New Cookie ", UUID)
		} else if UUIDSignCookie == nil {
			http.Error(rw, "No UUID sign", http.StatusUnauthorized)
			return
		} else {
			logrus.Debug("Check Cookie ", UUIDCookie.Value)
			h := hmac.New(sha256.New, []byte(key))
			h.Write([]byte(UUIDCookie.Value))
			referenceSign := h.Sum(nil)
			userSign, _ := hex.DecodeString(UUIDSignCookie.Value)
			if string(referenceSign) != string(userSign) {
				http.Error(rw, "Broken UUID sign", http.StatusUnauthorized)
			}
			UUID = UUIDCookie.Value
			isNewUUID = false
		}
		ctx := context.WithValue(r.Context(), RequestUUIDKey{}, UUID)
		ctx = context.WithValue(ctx, RequestisNewUUIDKey{}, isNewUUID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
