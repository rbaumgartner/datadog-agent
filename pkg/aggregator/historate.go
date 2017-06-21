package aggregator

// Historate tracks the distribution of samples added over one flush period for
// "rate" like metrics. Warning this doesn't use the harmonic mean, beware of
// what it means when using it.
type Historate struct {
	histogram         Histogram
	previousSample    float64
	previousTimestamp float64
	sampled           bool
}

func (h *Historate) addSample(sample *MetricSample, timestamp float64) {
	if h.previousTimestamp != 0 {
		v := (sample.Value - h.previousSample) / (timestamp - h.previousTimestamp)
		h.histogram.addSample(&MetricSample{Value: v}, timestamp)
		h.sampled = true
	}
	h.previousSample, h.previousTimestamp = sample.Value, timestamp
}

func (h *Historate) flush(timestamp float64) ([]*Serie, error) {
	if h.sampled == false {
		return []*Serie{}, NoSerieError{}
	}

	h.previousSample, h.previousTimestamp, h.sampled = 0.0, 0, false
	return h.histogram.flush(timestamp)
}