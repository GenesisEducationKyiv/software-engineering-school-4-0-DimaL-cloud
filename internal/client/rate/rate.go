package rate

type Rate interface {
	GetRate() (float64, error)
	GetNext() Rate
	SetNext(rateClient Rate) Rate
}
