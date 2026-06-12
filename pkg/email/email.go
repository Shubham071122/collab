package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shubham071122/collab/internal/config"
)

type SendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func SendVerificationEmail(toEmail string, name string, otp string) error {
	cfg := config.LoadConfig()
	apiKey := cfg.ResendAPIKey
	fromEmail := cfg.ResendFromEmail

	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not configured")
	}
	if fromEmail == "" {
		fromEmail = "no-reply@plynk.in"
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Verify Your Email</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #ffffff; margin: 0; padding: 0; -webkit-text-size-adjust: none;">
  <table width="100%%" cellpadding="0" cellspacing="0" role="presentation" style="background-color: #ffffff; width: 100%%; margin: 0; padding: 0;">
    <tr>
      <td align="center">
        <table width="100%%" cellpadding="0" cellspacing="0" role="presentation" style="max-width: 570px; margin: 0; padding: 0;">
          <!-- Email Body -->
          <tr>
            <td style="padding: 60px 20px; text-align: center;">
              <span style="font-size: 28px; font-weight: 800; color: #000; text-decoration: none; letter-spacing: -0.02em;">
                Collab
              </span>
            </td>
          </tr>
          <tr>
            <td style="padding: 0 40px;">
              <h1 style="font-size: 24px; font-weight: 700; color: #111827; margin-top: 0; text-align: center; letter-spacing: -0.02em;">
                Verify your account
              </h1>
              <p style="font-size: 16px; line-height: 24px; color: #4b5563; margin-bottom: 40px; text-align: center;">
                Hello %s, welcome to Collab. Use the code below to securely verify your email address.
              </p>
              
              <!-- Code Box -->
              <table width="100%%" border="0" cellspacing="0" cellpadding="0" role="presentation">
                <tr>
                  <td align="center" style="background-color: #f9fafb; padding: 30px; border-radius: 16px; border: 1px solid #f3f4f6;">
                    <span style="font-size: 42px; font-weight: 800; color: #000; letter-spacing: 12px; margin-left: 12px;">
                      %s
                    </span>
                  </td>
                </tr>
              </table>

              <p style="font-size: 14px; line-height: 20px; color: #6b7280; margin-top: 40px; text-align: center;">
                This code is valid for <strong>30 minutes</strong>. 
                If you didn't create an account, you can safely ignore this email.
              </p>
              
              <hr style="border: none; border-top: 1px solid #f3f4f6; margin: 40px 0;">
              
              <p style="font-size: 12px; line-height: 18px; color: #9ca3af; text-align: center;">
                &copy; %d Collab. All rights reserved.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, name, otp, time.Now().Year())

	reqBody := SendEmailRequest{
		From:    fmt.Sprintf("Collab <%s>", fromEmail),
		To:      []string{toEmail},
		Subject: "Your Collab verification code",
		HTML:    htmlContent,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend returned non-ok status: %s", resp.Status)
	}

	return nil
}
