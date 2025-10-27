package verifier

import "net/http"

type CustomVerifier struct {
	Secret string
}

func (v *CustomVerifier) Verify(r *http.Request, _ []byte) bool {
	providerSecret := r.Header.Get("X-Auth")

	return providerSecret != "" && providerSecret == v.Secret
}
