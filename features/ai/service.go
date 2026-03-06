package ai

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/TFX0019/api-go-gds/features/wallets"
	"github.com/TFX0019/api-go-gds/pkg/config"
)

type Service interface {
	CreateGeneration(userID uint, prompt string, imageInput *string) (*AIGenerationResponse, error)
	UpdateGenerationResult(id string, userID uint, imageOutput string, responseText string) (*AIGenerationResponse, error)
	GetGenerationsByUserID(userID uint, page, limit int) (*PaginatedAIGenerationResponse, error)
	GetAllGenerationsAdmin(page, limit int) (*PaginatedAIGenerationResponse, error)

	// AISuggestions
	CreateSuggestion(req CreateAISuggestionRequest) (*AISuggestionResponse, error)
	UpdateSuggestion(id string, req UpdateAISuggestionRequest) (*AISuggestionResponse, error)
	DeleteSuggestion(id string) error
	GetAllSuggestions() ([]AISuggestionResponse, error)
}

type service struct {
	repo        Repository
	walletsRepo wallets.Repository
}

func NewService(repo Repository, walletsRepo wallets.Repository) Service {
	return &service{repo: repo, walletsRepo: walletsRepo}
}

func (s *service) CreateGeneration(userID uint, prompt string, imageInput *string) (*AIGenerationResponse, error) {
	generation := &AIGeneration{
		UserID:     userID,
		Prompt:     prompt,
		ImageInput: imageInput,
	}

	if err := s.repo.Create(generation); err != nil {
		return nil, err
	}

	// Reload to get User and other preloads
	updatedGen, err := s.repo.FindByID(generation.ID.String())
	if err == nil {
		generation = updatedGen
	}

	res := mapToResponse(*generation)
	return &res, nil
}

func (s *service) UpdateGenerationResult(id string, userID uint, imageOutput string, responseText string) (*AIGenerationResponse, error) {
	generation, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("generation not found")
	}

	if generation.UserID != userID {
		return nil, errors.New("unauthorized to update this generation")
	}

	if generation.ImageOutput != nil || generation.ResponseText != nil {
		return nil, errors.New("this generation already has a result")
	}

	if imageOutput != "" {
		generation.ImageOutput = &imageOutput
	}
	if responseText != "" {
		generation.ResponseText = &responseText
	}
	if err := s.repo.Update(generation); err != nil {
		return nil, err
	}

	// Credits logic
	removeCreditsStr := config.GetEnv("REMOVE_CREDITS_FOR_GENERATION", "10")
	removeCredits, err := strconv.Atoi(removeCreditsStr)
	if err != nil {
		removeCredits = 10
	}

	refID := generation.ID.String()
	// Subtract credits (negative amount)
	if err := s.walletsRepo.AddCredits(userID, -removeCredits, wallets.TransactionTypeImageGeneration, &refID); err != nil {
		return nil, fmt.Errorf("failed to subtract credits: %v", err)
	}

	res := mapToResponse(*generation)
	return &res, nil
}

func (s *service) GetGenerationsByUserID(userID uint, page, limit int) (*PaginatedAIGenerationResponse, error) {
	offset := (page - 1) * limit
	gens, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var data []AIGenerationResponse
	for _, g := range gens {
		data = append(data, mapToResponse(g))
	}

	return &PaginatedAIGenerationResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetAllGenerationsAdmin(page, limit int) (*PaginatedAIGenerationResponse, error) {
	offset := (page - 1) * limit
	gens, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var data []AIGenerationResponse
	for _, g := range gens {
		data = append(data, mapToResponse(g))
	}

	return &PaginatedAIGenerationResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func mapToResponse(g AIGeneration) AIGenerationResponse {
	return AIGenerationResponse{
		ID:           g.ID.String(),
		UserID:       g.UserID,
		UserName:     g.User.Name,
		Prompt:       g.Prompt,
		ImageInput:   g.ImageInput,
		ImageOutput:  g.ImageOutput,
		ResponseText: g.ResponseText,
		CreatedAt:    g.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    g.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// AISuggestions Logic

func (s *service) CreateSuggestion(req CreateAISuggestionRequest) (*AISuggestionResponse, error) {
	suggestion := &AISuggestion{
		Prompt:      req.Prompt,
		Description: req.Description,
	}

	if err := s.repo.CreateSuggestion(suggestion); err != nil {
		return nil, err
	}

	res := mapSuggestionToResponse(*suggestion)
	return &res, nil
}

func (s *service) UpdateSuggestion(id string, req UpdateAISuggestionRequest) (*AISuggestionResponse, error) {
	suggestion, err := s.repo.FindSuggestionByID(id)
	if err != nil {
		return nil, errors.New("suggestion not found")
	}

	if req.Prompt != "" {
		suggestion.Prompt = req.Prompt
	}
	if req.Description != "" {
		suggestion.Description = req.Description
	}

	if err := s.repo.UpdateSuggestion(suggestion); err != nil {
		return nil, err
	}

	res := mapSuggestionToResponse(*suggestion)
	return &res, nil
}

func (s *service) DeleteSuggestion(id string) error {
	_, err := s.repo.FindSuggestionByID(id)
	if err != nil {
		return errors.New("suggestion not found")
	}
	return s.repo.DeleteSuggestion(id)
}

func (s *service) GetAllSuggestions() ([]AISuggestionResponse, error) {
	suggestions, err := s.repo.FindAllSuggestions()
	if err != nil {
		return nil, err
	}

	var res []AISuggestionResponse
	for _, sub := range suggestions {
		res = append(res, mapSuggestionToResponse(sub))
	}
	return res, nil
}

func mapSuggestionToResponse(s AISuggestion) AISuggestionResponse {
	return AISuggestionResponse{
		ID:          s.ID.String(),
		Prompt:      s.Prompt,
		Description: s.Description,
		CreatedAt:   s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
