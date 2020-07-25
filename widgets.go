package main

import (
	"cjavellana.me/launchpad/agent/errors"
	"cjavellana.me/launchpad/agent/metrics"
	"cjavellana.me/launchpad/agent/servers/nginx"
	"fmt"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"time"
)

func NewLineChart(data chan metrics.Metrics, maxSeriesElements int, timeInterval time.Duration) *linechart.LineChart {

	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	)

	if err != nil {
		panic(err)
	}

	go func() {
		var inputs []float64
		xSeriesLabels := make(map[int]string)

		for {
			d := <-data
			inputs = append(inputs, d.Value)

			now := time.Now()
			for i := 0; i <= maxSeriesElements; i++ {
				t := now.Add(timeInterval * -1)
				xSeriesLabels[maxSeriesElements - i] = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())

				now = t
			}

			// We only keep the last
			if len(inputs) > maxSeriesElements {
				inputs = inputs[1:]
			}

			if err := lc.Series("first", inputs,
				linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
				linechart.SeriesXLabels(xSeriesLabels),
			); err != nil {
				panic(err)
			}
		}
	}()

	return lc
}

type render func(*text.Text)

func NewRollContentDisplay(fn render) *text.Text {
	display, err := text.New(text.RollContent(), text.WrapAtWords())
	errors.CheckFatal(err)

	go fn(display)

	return display
}

func NewNginxStatusWindow(ch chan nginx.StandardStatus) *text.Text  {
	return NewRollContentDisplay(
		func(display *text.Text) {
			for {
				nginxStatus := <- ch

				if message := messageFrom(nginxStatus); message != "" {
					err := display.Write(message)
					errors.CheckFatal(err)
				}
			}
		})
}

func messageFrom(n nginx.StandardStatus) string {
	message := n.RawData
	if n.ErrorMessage != "" {
		message = n.ErrorMessage
	}

	return message
}

func NewServerStdoutWindow(ch chan string) *text.Text {
	return NewRollContentDisplay(
		func(display *text.Text) {
			for {
				if err := display.Write(fmt.Sprintf("%s", <-ch)); err != nil {
					panic(err)
				}
			}
		})
}
