package nginx

import (
	"bytes"
	"cjavellana.me/launchpad/agent/internal/pkg/cfg"
	"html/template"
)

type ConfigBuilder struct {
	config cfg.ServerConfig
}

func newConfigBuilder(c cfg.ServerConfig) *ConfigBuilder {
	return &ConfigBuilder{
		config: c,
	}
}

// Returns an nginx configuration or an error if
// there is an error building the configuration text
func (c *ConfigBuilder) Build() (string, error) {
	server := c.config

	tmpl, err := template.New("nginxCfg").Parse(server.ConfigTemplate)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, server)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
