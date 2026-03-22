package auth

import (
	"net/http"
	"time"
)

const SessionCookieName = "voidshift.session"

func BuildSessionCookie(sessionID string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
	}
}

func BuildExpiredSessionCookie() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}

}
