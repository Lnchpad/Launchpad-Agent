package view

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
)

type SimpleDashboard struct {
	cpuUsage    *linechart.LineChart
	memoryUsage *linechart.LineChart
	serverLog   *text.Text
	agentLog    *text.Text
}

func SimpleDashboardBuilder() *SimpleDashboard {
	return &SimpleDashboard{}
}

func (d *SimpleDashboard) WithCpuWidget(chart *linechart.LineChart) *SimpleDashboard {
	d.cpuUsage = chart
	return d
}

func (d *SimpleDashboard) WithMemoryWidget(chart *linechart.LineChart) *SimpleDashboard {
	d.memoryUsage = chart
	return d
}

func (d *SimpleDashboard) WithStdoutWidget(stdout *text.Text) *SimpleDashboard {
	d.serverLog = stdout
	return d
}

func (d *SimpleDashboard) WithNginxMetrics(nginxMetrics *text.Text) *SimpleDashboard {
	d.agentLog = nginxMetrics
	return d
}

func (d *SimpleDashboard) Build(terminal *termbox.Terminal) *container.Container {
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
