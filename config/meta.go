package config

var (
	MetaGraphBaseURL string
	MetaGraphVersion string
	MetaAccessToken  string
	MetaAdAccountID  string
)

func InitMeta() {
	MetaGraphBaseURL = GetEnv("META_GRAPH_BASE_URL", "https://graph.facebook.com")
	MetaGraphVersion = GetEnv("META_GRAPH_VERSION", "v25.0")
	MetaAccessToken = GetEnv("META_ACCESS_TOKEN", "")
	MetaAdAccountID = GetEnv("META_AD_ACCOUNT_ID", "act_772782991828692")
}
