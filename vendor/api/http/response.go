package http

type Response struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

func NewErrorData(ctx, error string) map[string]string {
	return map[string]string{"ctx":ctx, "error": error}
}
