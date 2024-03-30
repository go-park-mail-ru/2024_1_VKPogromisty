package middleware

import "net/http"

var (
	ALLOWED_HEADERS = []string{"Accept", "Accept-Language", "Content-Type"}
	ALLOWED_ORIGINS = []string{
		"http://localhost",
		"http://localhost:3000",
		"http://94.139.246.134",
		"http://socio-project.ru",
		"http://127.0.0.1",
		"http://127.0.0.1:3000",
	}
	ALLOWED_METHODS = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}
)
