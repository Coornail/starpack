package sampler

type ImageSampler interface {
	HasMore() bool
	Next() (x, y int)
	Reset()
}
