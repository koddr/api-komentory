package models

// PostmarkSuppressSendingWebhook struct to describe Postmark suppress sending webhook object.
//  - Recipient == subscriber email address;
//  - SuppressSending == true (deactivate) | false (reactivate);
// See: https://postmarkapp.com/developer/webhooks/subscription-change-webhook#subscription-change-webhook-data
type PostmarkSuppressSendingWebhook struct {
	Recipient       string `json:"Recipient" required:"required,email"`
	SuppressSending bool   `json:"SuppressSending"`
}
