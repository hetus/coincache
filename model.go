package coincache

type Model struct {
	Ask    float64
	Bid    float64
	High   float64
	ID     int `storm:"id,increment"`
	Last   float64
	Low    float64
	Name   string
	Volume float64
}
