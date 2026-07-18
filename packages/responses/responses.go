package response

import (
	"encoding/json"
	"idiom-api-services/domains/projects"
	"net/http"
	"net/url"
	"strings"
)

func JSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"error": message,
	})
}

func Message(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{
		"message": message,
	})
}

func Redirect(w http.ResponseWriter, r *http.Request, redirectURLs []projects.RedirectURL, code string) bool {
	if code == "" {
		return false
	}

	redirectURL := redirectURLForUserAgent(redirectURLs, r.UserAgent())
	if redirectURL == "" {
		return false
	}

	http.Redirect(w, r, redirectURLWithCode(redirectURL, code), http.StatusSeeOther)
	return true
}

func redirectURLForUserAgent(redirectURLs []projects.RedirectURL, userAgent string) string {
	device := deviceFromUserAgent(userAgent)
	for _, redirectURL := range redirectURLs {
		if strings.EqualFold(redirectURL.Device, device) && redirectURL.URL != "" {
			return redirectURL.URL
		}
	}
	for _, redirectURL := range redirectURLs {
		if strings.EqualFold(redirectURL.Device, "web") && redirectURL.URL != "" {
			return redirectURL.URL
		}
	}
	for _, redirectURL := range redirectURLs {
		if redirectURL.URL != "" {
			return redirectURL.URL
		}
	}
	return ""
}

func deviceFromUserAgent(userAgent string) string {
	value := strings.ToLower(userAgent)
	if strings.Contains(value, "iphone") || strings.Contains(value, "ipad") || strings.Contains(value, "android") || strings.Contains(value, "mobile") {
		return "mobile"
	}
	return "web"
}

func redirectURLWithCode(redirectURL, code string) string {
	parsed, err := url.Parse(redirectURL)
	if err != nil {
		return redirectURL
	}
	query := parsed.Query()
	query.Set("code", code)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
