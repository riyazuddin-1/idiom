package response

import (
	"encoding/json"
	"idiom-api-services/domains/projects"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta"`
}

type Meta struct {
	Timestamp  string      `json:"timestamp"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total   int `json:"total"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
	HasMore bool `json:"hasMore"`
}

func writeJSON(w http.ResponseWriter, status int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func OK(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func Created(w http.ResponseWriter, message string, data any) {
	writeJSON(w, http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func OKWithPagination(w http.ResponseWriter, message string, data any, pagination Pagination) {
	writeJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Pagination: &pagination,
		},
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, Response{
		Success: false,
		Message: message,
		Data:    nil,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
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
