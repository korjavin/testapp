package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	TelegramInitDataHeader = "X-Telegram-Init-Data"
)

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Language  string `json:"language_code,omitempty"`
}

type TelegramAuthResult struct {
	User   TelegramUser
	Valid  bool
	Errors []string
}

// ValidateInitData validates Telegram WebApp initData HMAC-SHA256 signature
func ValidateInitData(initData string) (*TelegramAuthResult, error) {
	result := &TelegramAuthResult{Valid: false}

	if initData == "" {
		result.Errors = append(result.Errors, "empty initData")
		return result, nil
	}

	// Parse initData
	values := make(map[string]string)
	for _, pair := range strings.Split(initData, "&") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			values[parts[0]] = parts[1]
		}
	}

	// Check hash
	hash := values["hash"]
	delete(values, "hash")

	// Check auth_date
	authDate, ok := values["auth_date"]
	if !ok {
		result.Errors = append(result.Errors, "missing auth_date")
		return result, nil
	}

	// Verify auth_date is not too old (24 hours max)
	t, err := time.Parse("2006-01-02T15:04:05Z", authDate)
	if err != nil {
		// Try Unix timestamp
		var sec int64
		if _, err := fmt.Sscanf(authDate, "%d", &sec); err == nil {
			t = time.Unix(sec, 0)
		} else {
			result.Errors = append(result.Errors, "invalid auth_date format")
			return result, nil
		}
	}

	if time.Since(t) > 24*time.Hour {
		result.Errors = append(result.Errors, "auth_date too old")
		return result, nil
	}

	// Build data check string
	var dataCheck []string
	for k, v := range values {
		dataCheck = append(dataCheck, fmt.Sprintf("%s=%s", k, v))
	}
	sortStrings(dataCheck)
	dataCheckString := strings.Join(dataCheck, "\n")

	// Verify HMAC
	botToken := os.Getenv("TG_BOT_TOKEN")
	secretKey := hmac.New(sha256.New, []byte("WebAppData")).Write([]byte(botToken)).Sum(nil)

	computedHash := hex.EncodeToString(hmac.New(sha256.New, secretKey).Sum([]byte(dataCheckString)))

	if computedHash != hash {
		result.Errors = append(result.Errors, "invalid hash")
		return result, nil
	}

	// Parse user
	var user TelegramUser
	if userData, ok := values["user"]; ok {
		// Simple parsing of user JSON
		// In production, use json.Unmarshal
		user = parseUserFromString(userData)
	}

	result.User = user
	result.Valid = true
	return result, nil
}

func parseUserFromString(data string) TelegramUser {
	// Simplified parser - in production use json.Unmarshal
	var user TelegramUser
	// This is a placeholder - proper implementation would use json.Unmarshal
	return user
}

func sortStrings(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func TelegramAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for certain paths
		if r.URL.Path == "/health" || r.URL.Path == "/static/" {
			next.ServeHTTP(w, r)
			return
		}

		initData := r.Header.Get(TelegramInitDataHeader)
		if initData == "" {
			// For development, allow访问 without auth
			if os.Getenv("APP_ENV") == "development" {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		result, err := ValidateInitData(initData)
		if err != nil {
			http.Error(w, "Auth error", http.StatusInternalServerError)
			return
		}

		if !result.Valid {
			http.Error(w, "Invalid auth", http.StatusUnauthorized)
			return
		}

		// Store user in context
		ctx := r.Context()
		ctx = WithTelegramUser(ctx, result.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type contextKey string

const telegramUserKey contextKey = "telegram_user"

func WithTelegramUser(ctx interface{}, user TelegramUser) interface{} {
	return contextWithValue(ctx, telegramUserKey, user)
}

func GetTelegramUser(ctx interface{}) (TelegramUser, bool) {
	val := contextValue(ctx, telegramUserKey)
	if val == nil {
		return TelegramUser{}, false
	}
	user, ok := val.(TelegramUser)
	return user, ok
}

func contextWithValue(ctx interface{}, key contextKey, value interface{}) interface{} {
	// This is simplified - proper implementation uses context.WithValue
	return ctx
}

func contextValue(ctx interface{}, key contextKey) interface{} {
	// This is simplified - proper implementation uses context.Value
	return nil
}
