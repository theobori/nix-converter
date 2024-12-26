package config

import "strings"

type ContainerConfig struct {
	Begin        string
	End          string
	Sep          string
	ElementBegin string
	ElementEnd   string
}

type Config struct {
	Map   ContainerConfig
	Array ContainerConfig

	IndentSize int
}

func (c *Config) IndentLevel() string {
	return strings.Repeat(" ", c.IndentSize)
}
