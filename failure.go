// Package failure implements the Phi Accrual Failure Detector
/*

http://ddg.jaist.ac.jp/pub/HDY+04.pdf

*/
package failure

import (
	"math"
	"time"

	"github.com/dgryski/go-onlinestats"
)

type Detector struct {
	w          *onlinestats.Windowed
	last       time.Time
	minSamples int
}

func New(windowSize, minSamples int) *Detector {

	d := &Detector{
		w:          onlinestats.NewWindowed(windowSize),
		minSamples: minSamples,
	}

	return d
}

func (d *Detector) Ping(now time.Time) {
	if !d.last.IsZero() {
		d.w.Push(now.Sub(d.last).Seconds())
	}
	d.last = now
}

func (d *Detector) Phi(now time.Time) float64 {
	if d.w.Len() < d.minSamples {
		return 0
	}

	t := now.Sub(d.last).Seconds()
	pLater := 1 - cdf(d.w.Mean(), d.w.Stddev(), t)
	phi := -math.Log10(pLater)

	return phi
}

func cdf(mean, stddev, x float64) float64 {
	return 0.5 + 0.5*math.Erf((x-mean)/(stddev*math.Sqrt2))
}
