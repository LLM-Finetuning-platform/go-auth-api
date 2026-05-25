package email

import (
	"github.com/eswarashish/go-auth-api/config"
	"github.com/resend/resend-go/v3"
)


type ResendClient struct {
	Createclient func (cfg *config.Config) (*resend.Client,error)
	Client *resend.Client
}

type EmailParams struct{
	From string
	To []string
	Html string
	Subject string
}

func NewClient  (cfg *config.Config) (*resend.Client, error){
	return resend.NewClient(cfg.ResendAPIKey), nil
}

func (rclient *ResendClient) EmailService (eparams *EmailParams) (*resend.SendEmailResponse, error){

		params := &resend.SendEmailRequest{
			From: eparams.From,
			To: eparams.To,
			Subject: eparams.Subject,
			Html: eparams.Html,
		}
		return rclient.Client.Emails.Send(params)

}