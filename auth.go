package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionCookieName = "session"
	sessionTTL        = 24 * time.Hour
	roleAdmin         = "admin"
	roleMint          = "mint"
)

func passHashAdmin() string {
	// h := os.Getenv("ADMIN_PASS_HASH")
	h := os.Getenv("ADMIN_PASS_HASH")

	if h == "" {
		log.Fatal("ADMIN_PASS_HASH is empty")
	}
	return h
}

func passHashMint() string {
	h := os.Getenv("MINT_PASS_HASH")
	if h == "" {
		log.Fatal("MINT_PASS_HASH is empty")
	}
	return h
}

func cookieSignKey() string {
	k := os.Getenv("COOKIE_SIGN_KEY")
	if len(k) < 32 {
		log.Fatal("COOKIE_SIGN_KEY is missing or too short (need 32+ chars)")
	}
	return k
}

func signSession(role string, expUnix int64) string {
	msg := role + "|" + strconv.FormatInt(expUnix, 10)
	m := hmac.New(sha256.New, []byte(cookieSignKey()))
	m.Write([]byte(msg))
	sig := hex.EncodeToString(m.Sum(nil))
	return msg + "|" + sig
}

func verifySession(token string) (string, bool) {
	parts := strings.Split(token, "|")
	if len(parts) != 3 {
		return "", false
	}

	role := parts[0]
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", false
	}
	if time.Now().Unix() > exp {
		return "", false
	}

	msg := parts[0] + "|" + parts[1]
	m := hmac.New(sha256.New, []byte(cookieSignKey()))
	m.Write([]byte(msg))
	want := hex.EncodeToString(m.Sum(nil))

	if !hmac.Equal([]byte(want), []byte(parts[2])) {
		return "", false
	}

	return role, true
}

func doLogin(c *gin.Context) {
	pass := c.PostForm("password")

	role := ""
	if bcrypt.CompareHashAndPassword([]byte(passHashAdmin()), []byte(pass)) == nil {
		role = roleAdmin
	} else if bcrypt.CompareHashAndPassword([]byte(passHashMint()), []byte(pass)) == nil {
		role = roleMint
	} else {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Wrong password"})
		return
	}

	exp := time.Now().Add(sessionTTL).Unix()
	token := signSession(role, exp) // в сессии храним роль

	secure := os.Getenv("ENV") != "dev"
	c.SetCookie(sessionCookieName, token, int(sessionTTL.Seconds()), "/", "", secure, true)

	// редирект сразу на нужную страницу
	if role == roleAdmin {
		c.Redirect(http.StatusFound, "/admin")
	} else {
		c.Redirect(http.StatusFound, "/mint")
	}
}

func doLogout(c *gin.Context) {
	secure := os.Getenv("ENV") != "dev"
	c.SetCookie(sessionCookieName, "", -1, "/", "", secure, true)
	c.Redirect(http.StatusFound, "/login")
}

func ginAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, err := c.Cookie(sessionCookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		role, ok := verifySession(tok)
		if !ok {
			clearSessionCookie(c)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("role", role)
		c.Next()
	}
}

func clearSessionCookie(c *gin.Context) {
	secure := os.Getenv("ENV") != "dev"
	c.SetCookie(sessionCookieName, "", -1, "/", "", secure, true)
}

func secureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	}
}