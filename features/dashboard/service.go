package dashboard

type Service interface {
	GetSummary() ([]SummaryItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSummary() ([]SummaryItem, error) {
	return s.repo.GetSummaryCounts()
}
