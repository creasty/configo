package configo

const (
	DEFAULT_CONFIG_ENV  = "development"
	DEFAULT_PREFIX      = "app"
	DEFAULT_DIR         = "."
	DEFAULT_DOTENV_PATH = ".env"
)

type Option struct {
	ConfigEnv  string
	Prefix     string
	Dir        string
	DotenvPath string
}
