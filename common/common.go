package common

type Question string

type RequestFormat struct {
	UserID    string     `json:"UserID"`
	Questions []Question `json:"Questions"`
}
