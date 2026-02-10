package utils

// set cookie options for http response
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetSecureCookie(c *gin.Context, token string, name string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    token,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false, // True on Prod
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Hoặc SameSiteStrictMode
	})
}

// clear cookie
func ClearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false, // True on Prod
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // Hoặc SameSiteStrictMode
	})
}
