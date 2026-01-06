package config

type Config struct {
	UploadPath              string `validate:"required"`
	DatabasePath            string `validate:"required"`
	MaxUploadSize           int64  `validate:"gt=0"`
	STTProvider             string `validate:"required,oneof=stub whisper"`
	ServerPort              int    `validate:"min=1,max=65535"`
	OAuthGoogleClientID     string `validate:"required"`
	OAuthGoogleClientSecret string `validate:"required"`
	AppBaseURL              string `validate:"required"`
	JWTSecret               string `validate:"required"`
	NLPServer               string `validate:"required"`
	APPUIBaseURL            string `validate:"required"`
}
