package indicator

import (
	"github.com/c9s/bbgo/pkg/types"
)

// Refer: Ease of Movement
// Refer URL: https://www.investopedia.com/terms/e/easeofmovement.asp

//go:generate callbackgen -type EMV
type EMV struct {
	types.IntervalWindow
	prevH    float64
	prevL    float64
	Values   *SMA
	EMVScale float64

	UpdateCallbacks []func(value float64)
}

const DefaultEMVScale float64 = 100000000.

func (inc *EMV) Update(high, low, vol float64) {
	if inc.EMVScale == 0 {
		inc.EMVScale = DefaultEMVScale
	}
	if inc.prevH == 0 || inc.Values == nil {
		inc.prevH = high
		inc.prevL = low
		inc.Values = &SMA{IntervalWindow: inc.IntervalWindow}
		return
	}
	distanceMoved := (high+low)/2. - (inc.prevH+inc.prevL)/2.
	boxRatio := vol / inc.EMVScale / (high - low)
	result := distanceMoved / boxRatio
	inc.prevH = high
	inc.prevL = low
	inc.Values.Update(result)
}

func (inc *EMV) Index(i int) float64 {
	if inc.Values == nil {
		return 0
	}
	return inc.Values.Index(i)
}

func (inc *EMV) Last() float64 {
	if inc.Values == nil {
		return 0
	}
	return inc.Values.Last()
}

func (inc *EMV) Length() int {
	if inc.Values == nil {
		return 0
	}
	return inc.Values.Length()
}

var _ types.Series = &EMV{}

func (inc *EMV) calculateAndUpdate(allKLines []types.KLine) {
	if inc.Values == nil {
		for _, k := range allKLines {
			inc.Update(k.High.Float64(), k.Low.Float64(), k.Volume.Float64())
			if inc.Length() > 0 {
				inc.EmitUpdate(inc.Last())
			}
		}
	} else {
		k := allKLines[len(allKLines)-1]
		inc.Update(k.High.Float64(), k.Low.Float64(), k.Volume.Float64())
		inc.EmitUpdate(inc.Last())
	}
}

func (inc *EMV) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.calculateAndUpdate(window)
}

func (inc *EMV) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
