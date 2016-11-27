package configo

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

	if self.ConfigEnv == "" {
		self.ConfigEnv = "development"
	}
}
