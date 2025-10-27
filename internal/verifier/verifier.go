package verifier

import "net/http"

type Verifier interface {
	Verify(r *http.Request, body []byte) bool
}
