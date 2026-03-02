package mocks

import "github.com/parkerjohnson/intercede/services/models"

// MockSubmissionService implements webhandlers.submissionLookup.
type MockSubmissionService struct {
	GetSubmissionsByChurchCodeFn func(churchCode string) (*[]models.SubmissionDB, error)
	CreatePrayerRequestFn        func(sub *models.Submission) error
	CreatePraiseReportFn         func(sub *models.Submission) error
	GetChurchCodesFn             func() ([]string, error)
}

func (m *MockSubmissionService) GetSubmissionsByChurchCode(churchCode string) (*[]models.SubmissionDB, error) {
	return m.GetSubmissionsByChurchCodeFn(churchCode)
}

func (m *MockSubmissionService) CreatePrayerRequest(sub *models.Submission) error {
	return m.CreatePrayerRequestFn(sub)
}

func (m *MockSubmissionService) CreatePraiseReport(sub *models.Submission) error {
	return m.CreatePraiseReportFn(sub)
}

func (m *MockSubmissionService) GetChurchCodes() ([]string, error) {
	return m.GetChurchCodesFn()
}
