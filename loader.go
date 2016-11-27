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
	self.Initialize()
	self.SetConfigEnv()

	if err := self.LoadFiles(); err != nil {
		return err
	}

	if err := envconfig.Process(self.Option.Prefix, self.Struct); err != nil {
		return err
	}

	if _, err := govalidator.ValidateStruct(self.Struct); err != nil {
		return err
	}

	return nil
}

func (self *Loader) Initialize() {
	opt := &self.Option
	opt.setDefault()

	godotenv.Load()
}

func (self *Loader) SetConfigEnv() {
	env := &struct {
		Env string
	}{}

	envconfig.Process(self.Option.Prefix, env)

	if env.Env != "" {
		self.Option.ConfigEnv = env.Env
	}
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
			if err := self.LoadFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *Loader) LoadFile(path string) error {
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
