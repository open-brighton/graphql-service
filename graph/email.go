package graph

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

const CONTACT_EMAIL = "josephcuffney@gmail.com"
const NO_REPLY_EMAIL = "no-reply@openbrighton.org"

func SendEmail(ctx context.Context, to, from, templateName, templateData string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	svc := ses.NewFromConfig(cfg)
	input := &ses.SendTemplatedEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Source:               aws.String(from),
		Template:             aws.String(templateName),
		TemplateData:         aws.String(templateData),
		ConfigurationSetName: aws.String("send-logs"), // enables logging

	}
	resp, err := svc.SendTemplatedEmail(ctx, input)
	if err != nil {
		log.Printf("Error sending email: %v", err)
	} else {
		pretty, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
		} else {
			log.Printf("SendTemplatedEmail response:\n%s", pretty)
		}
	}

	return err
}
