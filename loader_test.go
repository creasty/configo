package configo_test

import (
	"os"
	"testing"

	"github.com/creasty/configo"
)

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

func Test_Loader_PopulateOption(t *testing.T) {
	s := &Specification{}

	l := &configo.Loader{
		Struct: s,
	}

	{
		os.Clearenv()

		l.Option = configo.Option{}

		for actual, expect := range map[string]string{
			l.Option.Prefix:    "",
			l.Option.Dir:       "",
			l.Option.ConfigEnv: "",
		} {
			if actual != expect {
				t.Error("should be initialized")
				break
			}
		}

		l.PopulateOption()

		for actual, expect := range map[string]string{
			l.Option.Prefix:    configo.DEFAULT_PREFIX,
			l.Option.Dir:       configo.DEFAULT_DIR,
			l.Option.ConfigEnv: configo.DEFAULT_CONFIG_ENV,
		} {
			if actual != expect {
				t.Error("should set default values")
				break
			}
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

		for actual, expect := range map[string]string{
			l.Option.Prefix:    "prefix",
			l.Option.Dir:       "dir",
			l.Option.ConfigEnv: "configEnv",
		} {
			if actual != expect {
				t.Error("should not override with default values if values are set")
				break
			}
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

func Test_Loader_LoadFiles(t *testing.T) {
	l := &configo.Loader{
		Option: configo.Option{
			Dir: "./test/config",
		},
	}

	{
		os.Clearenv()

		s := &Specification{}
		l.Struct = s
		l.Option.ConfigEnv = "development"

		if err := l.LoadFiles(); err != nil {
			t.Error("should not fail", err)
		}

		for actual, expect := range map[[2]string]string{
			{s.Default, "Default"}:                                                                         "default",
			{s.Production, "Production"}:                                                                   "",
			{s.ProductionLocal, "ProductionLocal"}:                                                         "",
			{s.DefaultOverridedByProduction, "DefaultOverridedByProduction"}:                               "default",
			{s.DefaultOverridedByProductionLocal, "DefaultOverridedByProductionLocal"}:                     "default",
			{s.ProductionOverridedByProductionLocal, "ProductionOverridedByProductionLocal"}:               "",
			{s.Nested.Default, "Nested.Default"}:                                                           "default",
			{s.Nested.Production, "Nested.Production"}:                                                     "",
			{s.Nested.ProductionLocal, "Nested.ProductionLocal"}:                                           "",
			{s.Nested.DefaultOverridedByProduction, "Nested.DefaultOverridedByProduction"}:                 "default",
			{s.Nested.DefaultOverridedByProductionLocal, "Nested.DefaultOverridedByProductionLocal"}:       "default",
			{s.Nested.ProductionOverridedByProductionLocal, "Nested.ProductionOverridedByProductionLocal"}: "",
			{s.Subdir.Default, "Subdir.Default"}:                                                           "default",
			{s.Subdir.Production, "Subdir.Production"}:                                                     "",
			{s.Subdir.ProductionLocal, "Subdir.ProductionLocal"}:                                           "",
			{s.Subdir.DefaultOverridedByProduction, "Subdir.DefaultOverridedByProduction"}:                 "default",
			{s.Subdir.DefaultOverridedByProductionLocal, "Subdir.DefaultOverridedByProductionLocal"}:       "default",
			{s.Subdir.ProductionOverridedByProductionLocal, "Subdir.ProductionOverridedByProductionLocal"}: "",
		} {
			if actual[0] != expect {
				t.Errorf("expect %s to be %q, but was %q", actual[1], expect, actual[0])
			}
		}
	}

	{
		os.Clearenv()

		s := &Specification{}
		l.Struct = s
		l.Option.ConfigEnv = "production"

		if err := l.LoadFiles(); err != nil {
			t.Error("should not fail", err)
		}

		for actual, expect := range map[[2]string]string{
			{s.Default, "Default"}:                                                                         "default",
			{s.Production, "Production"}:                                                                   "production",
			{s.ProductionLocal, "ProductionLocal"}:                                                         "production.local",
			{s.DefaultOverridedByProduction, "DefaultOverridedByProduction"}:                               "production",
			{s.DefaultOverridedByProductionLocal, "DefaultOverridedByProductionLocal"}:                     "production.local",
			{s.ProductionOverridedByProductionLocal, "ProductionOverridedByProductionLocal"}:               "production.local",
			{s.Nested.Default, "Nested.Default"}:                                                           "default",
			{s.Nested.Production, "Nested.Production"}:                                                     "production",
			{s.Nested.ProductionLocal, "Nested.ProductionLocal"}:                                           "production.local",
			{s.Nested.DefaultOverridedByProduction, "Nested.DefaultOverridedByProduction"}:                 "production",
			{s.Nested.DefaultOverridedByProductionLocal, "Nested.DefaultOverridedByProductionLocal"}:       "production.local",
			{s.Nested.ProductionOverridedByProductionLocal, "Nested.ProductionOverridedByProductionLocal"}: "production.local",
			{s.Subdir.Default, "Subdir.Default"}:                                                           "default",
			{s.Subdir.Production, "Subdir.Production"}:                                                     "production",
			{s.Subdir.ProductionLocal, "Subdir.ProductionLocal"}:                                           "production.local",
			{s.Subdir.DefaultOverridedByProduction, "Subdir.DefaultOverridedByProduction"}:                 "production",
			{s.Subdir.DefaultOverridedByProductionLocal, "Subdir.DefaultOverridedByProductionLocal"}:       "production.local",
			{s.Subdir.ProductionOverridedByProductionLocal, "Subdir.ProductionOverridedByProductionLocal"}: "production.local",
		} {
			if actual[0] != expect {
				t.Errorf("expect %s to be %q, but was %q", actual[1], expect, actual[0])
			}
		}
	}
}

func Test_Loader_LoadEnvVars(t *testing.T) {
	l := &configo.Loader{
		Option: configo.Option{Prefix: "app"},
	}

	{
		os.Clearenv()

		s := &Specification{}
		l.Struct = s

		os.Setenv("APP_DEFAULT", "env")
		os.Setenv("APP_NESTED_DEFAULT", "env")

		if err := l.LoadEnvVars(); err != nil {
			t.Error("should not fail")
		}

		for actual, expect := range map[string]string{
			s.Default:        "env",
			s.Nested.Default: "env",
		} {
			if actual != expect {
				t.Error("should read values from env vars")
				break
			}
		}
	}

	{
		os.Clearenv()

		s := &Specification{}
		l.Struct = s

		if err := l.LoadEnvVars(); err != nil {
			t.Error("should not fail")
		}

		for actual, expect := range map[string]string{
			s.Production:        "dotenv",
			s.Nested.Production: "dotenv",
		} {
			if actual != expect {
				t.Error("should read values from dotenv file")
				break
			}
		}
	}
}

func Test_Loader_Validate(t *testing.T) {
	s := &Specification{}

	l := &configo.Loader{
		Struct: s,
		Option: configo.Option{},
	}

	if err := l.Validate(); err == nil {
		t.Error("should return error if a struct is not valid")
	}

	s.Required = "something"

	if err := l.Validate(); err != nil {
		t.Error("should validate a struct")
	}
}
