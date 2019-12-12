package types

type Answers struct {
	Items []Item `json:"items,omitempty`
}

type Item struct {
	Link string `json:"link,omitempty"`
}
