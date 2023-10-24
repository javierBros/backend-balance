package services

import "github.com/javierBros/backend-balance/application/model"

type ISummaryProcessingService interface {
	ProcessSummary(transactions []model.Transaction) error
}
