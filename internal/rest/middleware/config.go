package middleware

const (
	ALLOWED_HEADERS = "Accept, Accept-Language, Content-Type"
	ALLOWED_METHODS = "GET, POST, PUT, DELETE, OPTIONS"
)

var (
	ALLOWED_ORIGINS = []string{"http://localhost:3000", "http://94.139.246.134"}
)

func CheckAllowedOrigin(origin string) bool {
	for _, v := range ALLOWED_ORIGINS {
		if v == origin {
			return true
		}
	}
	return false
}
