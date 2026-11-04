package config

type Jobs struct {
	Default string         `yaml:"default"`
	Entries map[string]Job `yaml:",inline"`
}

type Job struct {
	Period string `yaml:"period"`
}
