package rate

type Rate interface {
	GetRate() (Response, error)
	GetNext() Rate
	SetNext(rateClient Rate) Rate
}

type Response struct {
	RateValue   float64
	APIResponse string
	APIName     string
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
