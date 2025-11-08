package httptransport

type VerificationData struct {
	Body        []byte
	Headers     map[string]string
	QueryParams map[string]string
}

type InboundAdapter interface {
	Verify(data VerificationData) bool
}
