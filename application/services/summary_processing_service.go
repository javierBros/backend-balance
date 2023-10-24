package services

type SummaryProcessingService struct {
	//emailSender adapters.IEmailSender
}

func (d *SummaryProcessingService) ProcessSummary() error {
	return nil
}

func NewSummaryProcessingService() *SummaryProcessingService {
	return &SummaryProcessingService{}
}
