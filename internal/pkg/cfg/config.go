package cfg

import (
	"cjavellana.me/launchpad/agent/internal/pkg/messaging/api"
	"flag"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type FsUpdaterConfig struct {
	// the location of all the static files in the
	// local file system i.e. in the container
	RootDirectory string

	// the url of the nexus repository
	NexusUrl string
}

type PortalApp struct {
	AppName string
}

type ServerConfig struct {
	ExecutablePath string
	ConfigTemplate string
	ConfigLocation string
	RootDirectory  string
	Applications   []PortalApp
}

var (
	defaultMaxSeriesElements = 15
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

type WebStatsConfig struct {
	StatsUrl         string
	SamplingInterval uint
	InitialDelay     uint
}

type ProbeConfig struct {
	Enabled          bool
	SamplingInterval uint
	WebStatsConfig   WebStatsConfig
	ProbeTypes       []ProbeType
}

type AppConfig struct {
	ViewType ViewType
	// the number of elements to retain in the charts
	// only applicable when ViewType is dashboard. When ViewType is none,
	// this property is ignored
	MaxSeriesElements int
	ServerConfig      ServerConfig
	ProbeConfig       ProbeConfig

	// messaging configurations
	BrokerConfig api.BrokerConfig `yaml:"brokerconfig,omitempty"`

	// The file system syncher configurations
	FsUpdaterConfig FsUpdaterConfig
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
	if appCfg.MaxSeriesElements == 0 {
		appCfg.MaxSeriesElements = defaultMaxSeriesElements
	}

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


