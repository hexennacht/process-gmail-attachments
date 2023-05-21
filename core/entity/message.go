package entity

type ReadFromRequest struct {
	EmailFrom    string `json:"email_from"`
	EmailSubject string `json:"email_subject"`
}
