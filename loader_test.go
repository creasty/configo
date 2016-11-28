package configo_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/creasty/configo"
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	dir = path.Join(dir, "test")
	os.Chdir(dir)
}

type Specification struct {
	Required string `valid:"required"`

	Default                              string
	Production                           string
	ProductionLocal                      string
	DefaultOverridedByProduction         string
	DefaultOverridedByProductionLocal    string
	ProductionOverridedByProductionLocal string

	Nested struct {
		Default                              string
		Production                           string
		ProductionLocal                      string
		DefaultOverridedByProduction         string
		DefaultOverridedByProductionLocal    string
		ProductionOverridedByProductionLocal string
	}

	Subdir struct {
		Default                              string
		Production                           string
		ProductionLocal                      string
		DefaultOverridedByProduction         string
		DefaultOverridedByProductionLocal    string
		ProductionOverridedByProductionLocal string
	}
}

var defaultOpt = configo.Option{
	Dir: "./config",
}

func Test_Loader_PopulateOption(t *testing.T) {
	s := &Specification{}

	l := &configo.Loader{
		Struct: s,
	}

	{
		os.Clearenv()

		l.Option = configo.Option{}

		if !(l.Option.Prefix == "" &&
			l.Option.Dir == "" &&
			l.Option.ConfigEnv == "") {
			t.Error("should be initialized")
		}

		l.PopulateOption()

		if !(l.Option.Prefix == configo.DEFAULT_PREFIX &&
			l.Option.Dir == configo.DEFAULT_DIR &&
			l.Option.ConfigEnv == configo.DEFAULT_CONFIG_ENV) {
			t.Error("should set default values")
		}
	}

	{
		os.Clearenv()

		l.Option = configo.Option{
			Prefix:    "prefix",
			Dir:       "dir",
			ConfigEnv: "configEnv",
		}

		l.PopulateOption()

		if !(l.Option.Prefix == "prefix" &&
			l.Option.Dir == "dir" &&
			l.Option.ConfigEnv == "configEnv") {
			t.Error("should not override with default values if values are set")
		}
	}

	{
		os.Clearenv()

		l.Option = configo.Option{}

		os.Setenv("APP_ENV", "production")
		l.PopulateOption()

		if l.Option.ConfigEnv != "production" {
			t.Error("should guess ConfigEnv from an env var")
		}
	}

	{
		os.Clearenv()

		l.Option = configo.Option{ConfigEnv: "configEnv"}

		os.Setenv("APP_ENV", "production")
		l.PopulateOption()

		if l.Option.ConfigEnv != "configEnv" {
			t.Error("should not override ConfigEnv with a value guessed from an env var if it's already set")
		}
	}
}

func Test_Loader_Validate(t *testing.T) {
	s := &Specification{}

	l := &configo.Loader{
		Struct: s,
		Option: defaultOpt,
	}

	if err := l.Validate(); err == nil {
		t.Error("should return error if a struct is not valid")
	}

	s.Required = "something"

	if err := l.Validate(); err != nil {
		t.Error("should validate a struct")
	}
}
