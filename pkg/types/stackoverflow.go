package types

type StackOverflowAnswers struct {
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	Link string `json:"link,omitempty"`
}
