package configo

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Option struct {
	ConfigEnv string
	Prefix    string
	Dir       string
}

func (self *Option) setDefault() {
	if self.Prefix == "" {
		self.Prefix = "app"
	}

	if self.Dir == "" {
		self.Dir = "."
	}
}

func Load(in interface{}, opt Option) error {
	_opt := &opt
	_opt.setDefault()

	godotenv.Load()

	if opt.ConfigEnv == "" {
		env := &struct {
			Env string
		}{}

		envconfig.Process(opt.Prefix, env)

		if env.Env != "" {
			opt.ConfigEnv = env.Env
		} else {
			opt.ConfigEnv = "development"
		}
	}

	{
		scopes := []string{
			"default",
			opt.ConfigEnv,
			opt.ConfigEnv + ".local",
		}
		for _, scope := range scopes {
			path := filepath.Join(opt.Dir, scope+".yml")
			if err := loadFile(in, path); err != nil {
				return err
			}
		}
	}

	if err := envconfig.Process(opt.Prefix, in); err != nil {
		return err
	}

	if _, err := govalidator.ValidateStruct(in); err != nil {
		return err
	}

	return nil
}

func loadFile(in interface{}, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, in)
}
