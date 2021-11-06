package timeseries

import (
	"github.com/K-Phoen/grabana/alert"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a graph panel.
type Option func(timeseries *TimeSeries)

// TooltipMode configures which series will be displayed in the tooltip
type TooltipMode string

const (
	// SingleSeries will only display the hovered series.
	SingleSeries TooltipMode = "single"
	// AllSeries will display all series.
	AllSeries = "multi"
	// NoSeries will hide the tooltip completely.
	NoSeries = "none"
)

// LegendOption allows to configure a legend.
type LegendOption uint16

const (
	// Hide keeps the legend from being displayed.
	Hide LegendOption = iota
	// AsTable displays the legend as a table.
	AsTable
	// AsList displays the legend as a list.
	AsList
	// Bottom displays the legend below the graph.
	Bottom
	// ToTheRight displays the legend on the right side of the graph.
	ToTheRight

	// Min displays the smallest value of the series.
	Min
	// Max displays the largest value of the series.
	Max
	// Avg displays the average of the series.
	Avg

	// First displays the first value of the series.
	First
	// FirstNonNull displays the first non-null value of the series.
	FirstNonNull
	// Last displays the last value of the series.
	Last
	// LastNonNull displays the last non-null value of the series.
	LastNonNull

	// Total displays the sum of values in the series.
	Total
	// Count displays the number of value in the series.
	Count
	// Range displays the difference between the minimum and maximum values.
	Range
)

// TimeSeries represents a time series panel.
type TimeSeries struct {
	Builder *sdk.Panel
}

// New creates a new time series panel.
func New(title string, options ...Option) *TimeSeries {
	panel := &TimeSeries{Builder: sdk.NewTimeseries(title)}
	panel.Builder.IsNew = false

	for _, opt := range append(defaults(), options...) {
		opt(panel)
	}

	return panel
}

func defaults() []Option {
	return []Option{
		Span(6),
		LineWidth(1),
		Tooltip(SingleSeries),
		Legend(Bottom, AsList),
	}
}

// DataSource sets the data source to be used by the graph.
func DataSource(source string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Datasource = &source
	}
}

// Tooltip configures the tooltip content.
func Tooltip(mode TooltipMode) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.Options.Tooltip.Mode = string(mode)
	}
}

// LineWidth defines the width of the line for a series (default 1, max 10, 0 is none).
func LineWidth(value int) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.TimeseriesPanel.FieldConfig.Defaults.Custom.LineWidth = value
	}
}

// Legend defines what should be shown in the legend.
func Legend(opts ...LegendOption) Option {
	return func(timeseries *TimeSeries) {
		legend := sdk.TimeseriesLegendOptions{
			DisplayMode: "list",
			Placement:   "bottom",
			Calcs:       make([]string, 0),
		}

		for _, opt := range opts {
			switch opt {
			case Hide:
				legend.DisplayMode = "hidden"
			case AsList:
				legend.DisplayMode = "list"
			case AsTable:
				legend.DisplayMode = "table"
			case ToTheRight:
				legend.Placement = "right"
			case Bottom:
				legend.Placement = "bottom"

			case First:
				legend.Calcs = append(legend.Calcs, "first")
			case FirstNonNull:
				legend.Calcs = append(legend.Calcs, "firstNotNull")
			case Last:
				legend.Calcs = append(legend.Calcs, "last")
			case LastNonNull:
				legend.Calcs = append(legend.Calcs, "lastNotNull")

			case Min:
				legend.Calcs = append(legend.Calcs, "min")
			case Max:
				legend.Calcs = append(legend.Calcs, "max")
			case Avg:
				legend.Calcs = append(legend.Calcs, "mean")

			case Count:
				legend.Calcs = append(legend.Calcs, "count")
			case Total:
				legend.Calcs = append(legend.Calcs, "sum")
			case Range:
				legend.Calcs = append(legend.Calcs, "range")
			}
		}

		timeseries.Builder.TimeseriesPanel.Options.Legend = legend
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Span = span
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Height = &height
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Description = &content
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Transparent = true
	}
}

// Alert creates an alert for this graph.
func Alert(name string, opts ...alert.Option) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Alert = alert.New(name, opts...).Builder
	}
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(timeseries *TimeSeries) {
		timeseries.Builder.Repeat = &repeat
	}
}
