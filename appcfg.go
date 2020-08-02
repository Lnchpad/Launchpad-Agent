package main

import (
	"bytes"
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/messaging"
	"cjavellana.me/launchpad/agent/servers/nginx"
	"github.com/go-yaml/yaml"
	"html/template"
	"io/ioutil"
)

type DashboardType string

const (
	simple DashboardType = "simple"
	none = "none"
)

type AppConfig struct {
	DashboardType DashboardType
	Server    nginx.Server
	Messaging messaging.KafkaConfig
}

// Reads the server config from the passed `configFile`.
func GetServerConfigFrom(configFile string) *AppConfig {
	appCfg := AppConfig{}

	cfgFile, err := ioutil.ReadFile(configFile)
	errors.CheckFatal(err)

	err = yaml.Unmarshal(cfgFile, &appCfg)
	errors.CheckFatal(err)

	server := appCfg.Server

	tmpl, err := template.New("nginxCfg").Parse(server.Config.Template)
	errors.CheckFatal(err)

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, server)
	errors.CheckFatal(err)

	server.Config.Template = tpl.String()

	return &appCfg
}
