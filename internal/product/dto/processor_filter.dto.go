package dto

type ProcessorFilterDTO struct {
	Brands      []string  `json:"brands" query:"brands"`           // ["Intel","Ryzen"]
	Frequencies []float64 `json:"frequencies" query:"frequencies"` // [2.5,2.6,2.8,...]
	Cores       []int     `json:"cores" query:"cores"`             // [1,2,4,6,8...]
	PriceAsc    bool      `json:"price_asc" query:"price_asc"`     // true = по возрастанию, false = по убыванию
	Limit       int       `json:"limit" query:"limit"`
	Offset      int       `json:"offset" query:"offset"`
}
