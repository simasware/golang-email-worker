package models

type SendEmailRequest struct {
	Recipient string `json:"recipient" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Body      string `json:"body" binding:"required"`
	Id        string `json:id`
}
