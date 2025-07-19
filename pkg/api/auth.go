package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/somepgs/go_final_project/tests"
	"net/http"
)

var secretKey = []byte("my_secret_key") // This should be a secure key, ideally loaded from an environment variable or a secure vault

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]any{"error": "Method not allowed"})
		return
	}

	if password == "" {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Password not set"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON format"})
		return
	}

	if req.Password != password {
		writeJson(w, http.StatusUnauthorized, map[string]any{"error": "Invalid password"})
		return
	}
	token, err := createJWT([]byte(password))
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
	writeJson(w, http.StatusOK, map[string]any{"token": token})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(password) > 0 {
			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			var valid bool
			if token != "" {
				valid, err = checkJWT(token, []byte(password))
				if err != nil {
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}
			if !valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
			tests.Token = token // Store the token in the tests package for testing purposes
		}
		next(w, r)
	})
}

func checkJWT(signedToken string, password []byte) (bool, error) {
	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secretKey, nil
	})
	if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, jwt.ErrTokenMalformed
	}

	storedHash, ok := claims["pwd_hash"].(string)
	if !ok {
		return false, jwt.ErrTokenMalformed
	}
	currentHash := sha256.Sum256(password)
	currentHashHex := hex.EncodeToString(currentHash[:])

	return storedHash == currentHashHex, nil
}

func createJWT(password []byte) (string, error) {
	hash := sha256.Sum256(password)
	hashHex := hex.EncodeToString(hash[:])

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pwd_hash": hashHex,
	})

	signedToken, err := tokenJWT.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
