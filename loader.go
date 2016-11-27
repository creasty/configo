package configo

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/mattn/go-zglob"
)

const FILE_GLOB_PATTERN = "%s/**/%s.yml"

type Loader struct {
	Option Option
	Struct interface{}
}

func (self *Loader) Load() error {
	self.PopulateOption()

	if err := self.LoadFiles(); err != nil {
		return err
	}

	if err := self.LoadEnvVars(); err != nil {
		return err
	}

	if err := self.Validate(); err != nil {
		return err
	}

	return nil
}

func (self *Loader) PopulateOption() {
	opt := &self.Option

	if opt.Prefix == "" {
		opt.Prefix = DEFAULT_PREFIX
	}

	if opt.Dir == "" {
		opt.Dir = DEFAULT_DIR
	}

	if opt.ConfigEnv == "" {
		env := &struct {
			Env string
		}{}

		envconfig.Process(opt.Prefix, env)

		if env.Env != "" {
			opt.ConfigEnv = env.Env
		} else {
			opt.ConfigEnv = DEFAULT_CONFIG_ENV
		}
	}
}

func (self *Loader) LoadEnvVars() error {
	godotenv.Load()
	return envconfig.Process(self.Option.Prefix, self.Struct)
}

func (self *Loader) LoadFiles() error {
	scopes := []string{
		"default",
		self.Option.ConfigEnv,
		self.Option.ConfigEnv + ".local",
	}

	for _, scope := range scopes {
		globPath := fmt.Sprintf(FILE_GLOB_PATTERN, self.Option.Dir, scope)
		paths, err := zglob.Glob(globPath)
		if err != nil {
			return err
		}

		for _, path := range paths {
			if err := self.loadFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *Loader) loadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, self.Struct)
}

func (self *Loader) Validate() error {
	_, err := govalidator.ValidateStruct(self.Struct)
	return err
}
