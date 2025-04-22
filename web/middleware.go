package web

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

func generateSessionToken() string {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// Base64 encode the random bytes
	return base64.URLEncoding.EncodeToString(b)
}

func setLoginCookie(w http.ResponseWriter, sessionToken string) {
	// Create a cookie with the session token
	cookie := &http.Cookie{
		Name:     "asbwig_session", // Name of the cookie
		Value:    sessionToken,    // Session token value
		Path:     "/",             // Make cookie available site-wide
		HttpOnly: true,            // Make cookie HTTP-only (can't be accessed via JavaScript)
		Expires:  time.Now().Add(24 * time.Hour), // Set expiration time for 24 hours
	}

	// Add cookie to response
	http.SetCookie(w, cookie)
}