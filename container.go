package main

type Container struct {
	cfg      *Config
	notebook *Notebook
}

func NewContainer(cfg *Config, path string) *Container {

	nb := NewNotebook(cfg, path)

	return &Container{
		cfg:      cfg,
		notebook: nb,
	}
}

func (s *Container) CurrentNotebook() (*Notebook, error) {
	return s.notebook, nil
}
