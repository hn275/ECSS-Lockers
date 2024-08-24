package admin

import (
	"net/http"

	"github.com/zvdv/ECSS-Lockers/internal/crypto"
	"github.com/zvdv/ECSS-Lockers/internal/env"
	"github.com/zvdv/ECSS-Lockers/internal/httputil"
	"github.com/zvdv/ECSS-Lockers/internal/logger"
)

var (
	adminUsername string
	adminPassword string
)

func Initialize() {
	adminUsername = env.MustEnv("ADMIN_USERNAME")
	adminPassword = env.MustEnv("ADMIN_PASSWORD")
}

func Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodGet {
		httputil.WriteResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if r.Method == http.MethodGet {
		httputil.WriteTemplatePage(w, nil, "templates/admin/auth.html")
		return
	}

	if err := r.ParseForm(); err != nil {
		logger.Error.Println(err)
		httputil.WriteResponse(w, http.StatusInternalServerError, nil)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username != adminUsername || password != adminPassword {
		httputil.WriteResponse(w, http.StatusOK, []byte(`
    <button type="submit" class="btn btn-primary btn-block">Login</button>
    <span id="error-message" class="text-error">Invalid credentials.</span>
            `))
		return
	}
	logger.Info.Println(username, password)
}

func makeToken() (string, error) {
	ciphertext, err := crypto.Encrypt(
		crypto.CipherKey[:],
		[]byte(adminPassword),
		[]byte(adminUsername))

	if err != nil {
		return "", err
	}

	token, err := crypto.SignMessage(
		crypto.SignatureKey[:],
		[]byte(adminPassword), ciphertext)

	if err != nil {
		return "", err
	}

	return crypto.Base64.EncodeToString(token), nil
}

func parseToken(token string) (bool, error) {
	tokenBytes, err := crypto.Base64.DecodeString(token)
	if err != nil {
		return false, nil
	}

	p := len(tokenBytes) - crypto.SignatureSize
	ciphertext, digest := tokenBytes[:p], tokenBytes[p:]

    pt, err := crypto.Decrypt(
        crypto.CipherKey[:],
        ciphertext,
        []byte(adminUsername))
    
    if err != nil {
        return false, err
    }

    return crypto.VerifySignature(crypto.SignatureKey[:], pt, digest)
}
