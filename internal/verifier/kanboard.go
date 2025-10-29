package verifier

import "net/http"

type KanboardVerifier struct {
	Secret string
}

func (v *KanboardVerifier) Verify(r *http.Request, _ []byte) bool {
	providerSecret := r.Header.Get("X-Kanboard-Token")

	return providerSecret != "" && providerSecret == v.Secret
}
