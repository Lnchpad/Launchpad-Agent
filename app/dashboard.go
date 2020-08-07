package app

import (
	"cjavellana.me/launchpad/agent/app/cfg"
	"cjavellana.me/launchpad/agent/app/view/widgets"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"time"
)

type SimpleDashboard struct {
	cpuUsage    *widgets.LineChart
	memoryUsage *widgets.LineChart
	serverLog   *widgets.RollContentDisplay
	agentLog    *widgets.RollContentDisplay
}

func NewSimpleDashboardBuilder(agent *Agent) *SimpleDashboard {
	// sync the widget time interval to that of the probe sampling intervals
	timeInterval := samplingInterval(agent)

	memWidget := widgets.NewLineChart(agent.AppCfg.SeriesElements, time.Duration(timeInterval)*time.Second)
	cpuWidget := widgets.NewLineChart(agent.AppCfg.SeriesElements, time.Duration(timeInterval)*time.Second)
	serverLog := widgets.NewRollContentDisplay()
	agentLog := widgets.NewRollContentDisplay()

	dashboard := SimpleDashboard{
		cpuUsage:    cpuWidget,
		memoryUsage: memWidget,
		serverLog:   serverLog,
		agentLog:    agentLog,
	}

	observeCpuAndMemProbes(&dashboard, agent)
	observeServerStdout(&dashboard, agent)

	return &dashboard
}

func observeServerStdout(s *SimpleDashboard, agent *Agent) {
	agent.Nginx.Process.Stdout.Observe(s.serverLog)
}

func observeCpuAndMemProbes(s *SimpleDashboard, agent *Agent) {
	for _, probe := range agent.Probes {
		switch probe.Type() {
		case cfg.CpuProbe:
			probe.Observe(s.cpuUsage)
		case cfg.MemProbe:
			probe.Observe(s.memoryUsage)
		}
	}
}

func samplingInterval(agent *Agent) uint {
	return agent.AppCfg.ProbeConfig.SamplingInterval
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