package models

type Config struct {
	Base        string `yaml:"base,omitempty"`
	Concurrency int    `yaml:"concurrency,omitempty"`
	Headers     []struct {
		Key   string `yaml:"key,omitempty"`
		Value string `yaml:"value,omitempty"`
	} `yaml:"headers,omitempty"`
	Paths      []string `yaml:"paths,omitempty"`
	StatusCode int      `yaml:"statuscode,omitempty"`
}
