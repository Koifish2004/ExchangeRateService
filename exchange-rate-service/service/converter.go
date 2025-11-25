package service

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) Convert(from, to string, amount float64, rates map[string]float64) float64 {

	fromRate := rates[from]

	toRate := rates[to]

	amountInUSD := amount / fromRate
	result := amountInUSD * toRate

	return result

}
