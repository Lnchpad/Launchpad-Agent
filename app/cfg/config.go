package cfg

import (
	"bytes"
	"cjavellana.me/launchpad/agent/app/messaging"
	"flag"
	"github.com/go-yaml/yaml"
	"html/template"
	"io/ioutil"
	"log"
)

type PortalApp struct {
}

type ServerConfig struct {
	ExecutablePath string
	ConfigTemplate string
	ConfigLocation string
	RootDirectory  string
	Applications   []PortalApp
}

const (
	DefaultMaxSeriesElements int = 15
)

type ViewType string

const (
	ViewTypeDashboardSimple ViewType = "dashboard-simple"
	ViewTypeNone                     = "none"
)

type ProbeType string

const (
	CpuProbe ProbeType = "cpu"
	MemProbe           = "memory"
)

type ProbeConfig struct {
	Enabled          bool
	SamplingInterval uint
	ProbeTypes       []ProbeType
}

type AppConfig struct {
	ViewType ViewType
	// the number of elements to retain in the charts
	// only applicable when ViewType is dashboard. When ViewType is none,
	// this property is ignored
	SeriesElements int
	ServerConfig   ServerConfig

	ProbeConfig ProbeConfig

	// messaging configurations
	BrokerConfig messaging.BrokerConfig `yaml:"brokerconfig,omitempty"`
}

func Get() AppConfig {
	cmdLineArgs := cmdLineArgs()
	return parseConfigYamlFile(cmdLineArgs.configFile)
}

func parseConfigYamlFile(configFile string) AppConfig {
	fileContentAsBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	return parseConfigYaml(fileContentAsBytes)
}

// Reads the server config from the passed `configFile`.
func parseConfigYaml(yamlConfig []byte) AppConfig {
	appCfg := AppConfig{}

	err := yaml.Unmarshal(yamlConfig, &appCfg)
	if err != nil {
		log.Fatal(err)
	}

	// set series element defaults (if none is given)
	if appCfg.SeriesElements == 0 {
		appCfg.SeriesElements = DefaultMaxSeriesElements
	}

	server := appCfg.ServerConfig

	tmpl, err := template.New("nginxCfg").Parse(server.ConfigTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, server)
	if err != nil {
		log.Fatal(err)
	}

	server.ConfigTemplate = tpl.String()

	return appCfg
}

type CmdLineArgs struct {
	configFile string
}

func cmdLineArgs() CmdLineArgs {
	var cfg string
	flag.StringVar(&cfg, "config", "config.yaml", "The configuration file to use")
	flag.Parse()

	return CmdLineArgs{configFile: cfg}
}
