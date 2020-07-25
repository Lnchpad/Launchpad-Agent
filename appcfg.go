package main

import (
	"bytes"
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/servers/nginx"
	"github.com/go-yaml/yaml"
	"html/template"
	"io/ioutil"
)

// Reads the server config from the passed `configFile`.
func GetServerConfigFrom(configFile string) *nginx.Server {
	server := nginx.Server{}

	cfgFile, err := ioutil.ReadFile(configFile)
	errors.CheckFatal(err)

	err = yaml.Unmarshal(cfgFile, &server)
	errors.CheckFatal(err)

	tmpl, err := template.New("nginxCfg").Parse(server.Config.Template)
	errors.CheckFatal(err)

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, server)
	errors.CheckFatal(err)

	server.Config.Template = tpl.String()

	return &server
}
