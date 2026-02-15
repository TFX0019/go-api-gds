package dashboard

type Service interface {
	GetSummary(userID string) ([]SummaryItem, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSummary(userID string) ([]SummaryItem, error) {
	return s.repo.GetSummaryCounts(userID)
}
