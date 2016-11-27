package configo

func Load(in interface{}, opt Option) error {
	l := &Loader{
		Struct: in,
		Option: opt,
	}
	return l.Load()
}
