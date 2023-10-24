package notifications

type IEmailSender interface {
	SendSummaryEmail(summary string) error
}
