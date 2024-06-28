package rate

type Rate interface {
	GetRate() (float64, error)
	GetNext() Rate
	SetNext(rateClient Rate) Rate
}

type ClientChain struct {
	next Rate
}

func (rcc *ClientChain) GetNext() Rate {
	return rcc.next
}

func (rcc *ClientChain) SetNext(next Rate) Rate {
	rcc.next = next
	return next
}
