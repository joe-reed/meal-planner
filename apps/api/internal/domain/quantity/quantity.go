package quantity

type Quantity struct {
	Amount int  `json:"amount"`
	Unit   Unit `json:"unit"`
}
