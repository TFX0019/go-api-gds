package plans

type Service interface {
	ListAll() ([]Plan, error)
	ListActive() ([]Plan, error)
	Update(id uint, dto PlanUpdateDTO) (*Plan, error)
	Activate(id uint) (*Plan, error)
	Deactivate(id uint) (*Plan, error)
	Create(dto PlanCreateDTO) (*Plan, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListAll() ([]Plan, error) {
	return s.repo.FindAll()
}

func (s *service) ListActive() ([]Plan, error) {
	return s.repo.FindAllActive()
}

func (s *service) Update(id uint, dto PlanUpdateDTO) (*Plan, error) {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Title != nil {
		plan.Title = *dto.Title
	}
	if dto.Description != nil {
		plan.Description = *dto.Description
	}
	if dto.Price != nil {
		plan.Price = *dto.Price
	}
	if dto.Benefits != nil {
		// Ensure benefits is not nil or empty if provided
		plan.Benefits = dto.Benefits
	}
	if dto.MaxCustomers != nil {
		plan.MaxCustomers = *dto.MaxCustomers
	}
	if dto.MaxProducts != nil {
		plan.MaxProducts = *dto.MaxProducts
	}
	if dto.MaxMaterials != nil {
		plan.MaxMaterials = *dto.MaxMaterials
	}
	if dto.MaxTasks != nil {
		plan.MaxTasks = *dto.MaxTasks
	}
	if dto.IsActive != nil {
		plan.IsActive = *dto.IsActive
	}

	if err := s.repo.Update(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *service) Activate(id uint) (*Plan, error) {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if plan.IsActive {
		return plan, nil // Already active
	}
	plan.IsActive = true
	if err := s.repo.Update(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *service) Deactivate(id uint) (*Plan, error) {
	plan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if !plan.IsActive {
		return plan, nil // Already inactive
	}
	plan.IsActive = false
	if err := s.repo.Update(plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *service) Create(dto PlanCreateDTO) (*Plan, error) {
	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}

	plan := &Plan{
		ProductID:    dto.ProductID,
		Title:        dto.Title,
		Description:  dto.Description,
		Price:        dto.Price,
		Benefits:     dto.Benefits,
		MaxCustomers: *dto.MaxCustomers,
		MaxProducts:  *dto.MaxProducts,
		MaxMaterials: *dto.MaxMaterials,
		MaxTasks:     *dto.MaxTasks,
		IsActive:     isActive,
	}

	if err := s.repo.Create(plan); err != nil {
		return nil, err
	}
	return plan, nil
}
