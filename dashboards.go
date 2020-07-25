package main

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

type StandardDashboard struct {
	terminal    *termbox.Terminal
	cpuUsage    *linechart.LineChart
	memoryUsage *linechart.LineChart
	serverLog   *text.Text
	agentLog    *text.Text
}

func StandardDashboardBuilder(terminal *termbox.Terminal) *StandardDashboard {
	return &StandardDashboard{terminal: terminal}
}

func (d *StandardDashboard) WithCpuWidget(chart *linechart.LineChart) *StandardDashboard {
	d.cpuUsage = chart
	return d
}

func (d *StandardDashboard) WithMemoryWidget(chart *linechart.LineChart) *StandardDashboard {
	d.memoryUsage = chart
	return d
}

func (d *StandardDashboard) WithStdoutWidget(stdout *text.Text) *StandardDashboard {
	d.serverLog = stdout
	return d
}

func (d *StandardDashboard) WithNginxMetrics(nginxMetrics *text.Text) *StandardDashboard {
	d.agentLog = nginxMetrics
	return d
}

func (d *StandardDashboard) build() *container.Container {
	cpuWidget := make([]container.Option, 0, 3)
	cpuWidget = append(cpuWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("CPU Usage"),
		container.PlaceWidget(d.cpuUsage),
	)

	memoryWidget := make([]container.Option, 0, 3)
	memoryWidget = append(memoryWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Memory Usage"),
		container.PlaceWidget(d.memoryUsage),
	)

	stdoutWidget := make([]container.Option, 0, 3)
	stdoutWidget = append(stdoutWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Server Logs"),
		container.PlaceWidget(d.serverLog),
	)

	agentWidget := make([]container.Option, 0, 3)
	agentWidget = append(agentWidget,
		container.Border(linestyle.Light),
		container.BorderTitle("Agent Logs"),
		container.PlaceWidget(d.agentLog),
	)

	c, err := container.New(
		d.terminal,
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
