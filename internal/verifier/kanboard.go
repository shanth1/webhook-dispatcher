package verifier

import (
	"net/http"
)

type KanboardVerifier struct {
	Secret string
}

func (v *KanboardVerifier) Verify(r *http.Request, _ []byte) bool {
	tokenFromRequest := r.URL.Query().Get("token")

	return tokenFromRequest != "" && tokenFromRequest == v.Secret
}
