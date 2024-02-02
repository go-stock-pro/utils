package token

const (
	ENV_AUTH_PUBLIC_KEY  = "AUTH_PUBLIC_KEY"
	ENV_AUTH_PRIVATE_KEY = "AUTH_PRIVATE_KEY"
)

type GenerateAuthTokenResponse struct {
	AuthToken string `json:"authToken"`
	ExpireAt  string `json:"expireAt"`
}
