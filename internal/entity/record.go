package entity

type Redirect struct {
	ActiveLink  string `json:"active_link,omitempty"`
	HistoryLink string `json:"history_link,omitempty"`
}
