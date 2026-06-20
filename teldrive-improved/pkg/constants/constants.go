package constants

const (
	AppName    = "TelDrive"
	AppVersion = "2.0.0"

	DefaultPageSize = 20
	MaxPageSize     = 100
	MaxFileSize     = 2 * 1024 * 1024 * 1024

	ChunkSize = 1024 * 1024

	JWTContextKey  = "jwt_claims"
	UserContextKey = "user"
)
