package models

type Company struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	ActiveUntil string `json:"active_until,omitempty"`
}
