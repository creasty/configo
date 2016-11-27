package configo

const (
	DEFAULT_CONFIG_ENV = "development"
	DEFAULT_PREFIX     = "app"
	DEFAULT_DIR        = "."
)

type Option struct {
	ConfigEnv string
	Prefix    string
	Dir       string
}
