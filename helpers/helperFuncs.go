package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GetQueryParameter returns value of a query parameter. Returns empty string if non-existing.
func GetQueryParameter(r *http.Request, name string) string {
	queryValue := r.URL.Query().Get(name)
	return queryValue
}

// Serialize converts an interface to a byte array.
func Serialize(input interface{}) ([]byte, error) {
	js, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	return js, nil
}

// IsFromSimulator checks if the request comes from the simulator.
func IsFromSimulator(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		data := map[string]interface{}{
			"status":    http.StatusForbidden,
			"error_msg": "You are not authorized to use this resource!",
		}

		jsonData, _ := Serialize(data)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(jsonData)
		return true
	}
	return false
}

// FormatDatetime formats a timestamp for display.
func FormatDatetime(timestamp int64) string {
	timeObject := time.Unix(timestamp, 0)
	return timeObject.Format("2006-02-01 @ 15:04")
}

// GetGravatarURL return the gravatar image for the given email address.
func GetGravatarURL(email string, size int) string {
	encodedEmail := hex.EncodeToString([]byte(strings.ToLower(strings.TrimSpace(email))))
	hashedEmail := fmt.Sprintf("%x", sha256.Sum256([]byte(encodedEmail)))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon&s=%d", hashedEmail, size)
}

var (
	_, _, _, _ = runtime.Caller(0)
)

// GetFullPath loads full path of file
func GetFullPath(fileName string) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unable to identify current directory (needed to load .env.test)")
		os.Exit(1)
	}
	basepath := filepath.Dir(file)
	return filepath.Join(basepath, fileName)
}
