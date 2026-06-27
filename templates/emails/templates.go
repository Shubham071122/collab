package emails

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed *.html
var templateFS embed.FS

var (
	subscriptionConfirmationTmpl *template.Template
)

func init() {
	var err error
	subscriptionConfirmationTmpl, err = template.ParseFS(templateFS, "subscription_confirmation.html")
	if err != nil {
		panic(err)
	}
}

type SubscriptionEmailData struct {
	Name            string
	PlanName        string
	Price           string
	SubscriptionID  string
	RenewalDate     string
	Features        []string
	LogoURL         string
	CTAURL          string
	Year            int
}

func RenderSubscriptionConfirmation(data SubscriptionEmailData) (string, error) {
	var buf bytes.Buffer
	err := subscriptionConfirmationTmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
