package dashboard

import (
	"cjavellana.me/launchpad/agent/internal/pkg/cfg"
	"cjavellana.me/launchpad/agent/internal/pkg/system"
	"cjavellana.me/launchpad/agent/internal/pkg/view/widgets"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"time"
)

type SimpleDashboardConfig struct {
	AppCfg    cfg.AppConfig
	Probes    []system.Probe
	ServerLog *system.Stdout
	WebStats  *system.WebStatsProbe
}

type SimpleDashboard struct {
	cpuUsage    *widgets.LineChart
	memoryUsage *widgets.LineChart
	serverLog   *widgets.RollContentDisplay
	agentLog    *widgets.RollContentDisplay
}

func NewSimpleDashboardBuilder(cfg SimpleDashboardConfig) *SimpleDashboard {
	// sync the widget time interval to that of the probe sampling intervals
	timeInterval := samplingInterval(cfg.AppCfg)

	memWidget := widgets.NewLineChart(cfg.AppCfg.MaxSeriesElements, time.Duration(timeInterval)*time.Second)
	cpuWidget := widgets.NewLineChart(cfg.AppCfg.MaxSeriesElements, time.Duration(timeInterval)*time.Second)
	serverLog := widgets.NewRollContentDisplay()
	agentLog := widgets.NewRollContentDisplay()

	dashboard := SimpleDashboard{
		cpuUsage:    cpuWidget,
		memoryUsage: memWidget,
		serverLog:   serverLog,
		agentLog:    agentLog,
	}

	observeCpuAndMemProbes(&dashboard, cfg.Probes)
	observeServerStdout(&dashboard, cfg.ServerLog)
	observeWebStats(&dashboard, cfg.WebStats)

	return &dashboard
}

func observeWebStats(s *SimpleDashboard, ws *system.WebStatsProbe) {
	ws.Observe(s.agentLog)
}

func observeServerStdout(s *SimpleDashboard, stdout *system.Stdout) {
	stdout.Observe(s.serverLog)
}

func observeCpuAndMemProbes(s *SimpleDashboard, probes []system.Probe) {
	for _, probe := range probes {
		switch probe.Type() {
		case cfg.CpuProbe:
			probe.Observe(s.cpuUsage)
		case cfg.MemProbe:
			probe.Observe(s.memoryUsage)
		}
	}
}

func samplingInterval(appCfg cfg.AppConfig) uint {
	return appCfg.ProbeConfig.SamplingInterval
}

func (d *SimpleDashboard) Build(terminal *termbox.Terminal) *container.Container {
	cpuWidget := make([]container.Option, 0, 3)
	cpuWidget = append(cpuWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("CPU Usage"),
		container.PlaceWidget(d.cpuUsage.LineChart),
	)

	memoryWidget := make([]container.Option, 0, 3)
	memoryWidget = append(memoryWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Memory Usage"),
		container.PlaceWidget(d.memoryUsage.LineChart),
	)

	stdoutWidget := make([]container.Option, 0, 3)
	stdoutWidget = append(stdoutWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Server Logs"),
		container.PlaceWidget(d.serverLog.Display),
	)

	agentWidget := make([]container.Option, 0, 3)
	agentWidget = append(agentWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Agent Logs"),
		container.PlaceWidget(d.agentLog.Display),
	)

	c, err := container.New(
		terminal,
		container.Border(linestyle.Light),
		container.BorderTitle("Launchpad WebServer Agent"),
		container.SplitVertical(
			container.Left(
				container.SplitHorizontal(
					container.Top(cpuWidget...),
					container.Bottom(memoryWidget...),
				),
			),
			container.Right(
				container.SplitHorizontal(
					container.Top(stdoutWidget...),
					container.Bottom(agentWidget...),
				),
			),
		),
	)

	if err != nil {
		panic(err)
	}

	return c
}
