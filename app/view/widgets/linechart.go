package widgets

import (
	"cjavellana.me/launchpad/agent/app/stats"
	"fmt"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
	"time"
)

type LineChart struct {
	LineChart *linechart.LineChart

	timeInterval time.Duration
	maxSeriesElements int

	inputs []float64
	xSeriesLabels map[int]string
}

func NewLineChart(maxSeriesElements int, timeInterval time.Duration) *LineChart {
	if lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorCyan)),
	); err != nil {
		panic(err)
	} else {
		return &LineChart{
			LineChart: lc,
			timeInterval: timeInterval,
			// +1 makes the line chart display maxSeriesElements time series
			maxSeriesElements: maxSeriesElements + 1,
			xSeriesLabels: make(map[int]string),
			inputs: make([]float64, 0, maxSeriesElements),
		}
	}
}

func (lc *LineChart) Update(m stats.Stats) {
	lc.inputs = append(lc.inputs, m.Value)

	now := time.Now()

	for i := 0; i <= lc.maxSeriesElements; i++ {
		t := now.Add(lc.timeInterval * -1)
		lc.xSeriesLabels[lc.maxSeriesElements - i] = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())

		now = t
	}

	// We only retain the last `maxSeriesElements`
	if len(lc.inputs) > lc.maxSeriesElements {
		lc.inputs = lc.inputs[1:]
		delete(lc.xSeriesLabels, lc.maxSeriesElements)
	}

	if err := lc.LineChart.Series("first", lc.inputs,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlue)),
		linechart.SeriesXLabels(lc.xSeriesLabels),
	); err != nil {
		panic(err)
	}
}
