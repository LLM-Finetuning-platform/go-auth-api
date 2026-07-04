package email

import (
	"fmt"

	"github.com/LLM-Finetuning-platform/go-auth-api/config"
	"github.com/go-playground/validator/v10"
	"github.com/resend/resend-go/v3"
)
type Request struct {
	Email  string `json:"email" validate:"required,email"`
	Client *ResendClient
	Params *EmailParams
	OTP    string
}

func ValidateEmailFormat(email string) error {
	validate := validator.New() 
	req := Request{Email: email}
	return validate.Struct(req)
}


func (req *Request) Verify() error {
	err := ValidateEmailFormat(req.Email)
	if err != nil {
		return fmt.Errorf("Improper Email Format %s", err)
	}

	return nil
}


type ResendClient struct {
	Createclient func(cfg *config.Config) (*resend.Client, error)
	Client       *resend.Client
}

type EmailParams struct {
	From    string
	To      []string
	Html    string
	Subject string
}

func NewClient(cfg *config.Config) (*resend.Client, error) {
	return resend.NewClient(cfg.ResendAPIKey), nil
}

func (rclient *ResendClient) EmailService(eparams *EmailParams) (*resend.SendEmailResponse, error) {

	params := &resend.SendEmailRequest{
		From:    eparams.From,
		To:      eparams.To,
		Subject: eparams.Subject,
		Html:    eparams.Html,
	}
	return rclient.Client.Emails.Send(params)

}
