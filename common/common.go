package common

type Question string
type UserID string

type RequestFormat struct {
	UserID    UserID     `json:"UserID"`
	Questions []Question `json:"Questions"`
}
