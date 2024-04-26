package interceptors

var (
	PublicMethods = map[string]struct{}{
		"/auth.Auth/Login":           {},
		"/auth.Auth/Logout":          {},
		"/auth.Auth/ValidateSession": {},
		"/user.User/GetByEmail":      {},
		"/user.User/Create":          {},
	}
)
