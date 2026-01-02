package services

import (
	"context"
	"strings"
	"time"

	"personal-budgeting/be/internal/clock"
	"personal-budgeting/be/internal/id"
	"personal-budgeting/be/internal/models"
	"personal-budgeting/be/internal/repositories"
)

type CategoryService struct {
	clk clock.Clock
	ids id.Generator

	cats repositories.CategoryRepository
}

func NewCategoryService(clk clock.Clock, ids id.Generator, cats repositories.CategoryRepository) *CategoryService {
	return &CategoryService{clk: clk, ids: ids, cats: cats}
}

func (s *CategoryService) List(ctx context.Context) ([]models.Category, error) {
	return s.cats.List(ctx)
}

func (s *CategoryService) Get(ctx context.Context, id string) (models.Category, error) {
	return s.cats.Get(ctx, id)
}

type CreateCategoryInput struct {
	Type        models.CategoryType `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
}

func (s *CategoryService) Create(ctx context.Context, in CreateCategoryInput) (models.Category, error) {
	now := s.clk.Now().Format(time.RFC3339)
	c := models.Category{
		ID:          s.ids.NewID(),
		Type:        in.Type,
		Name:        strings.TrimSpace(in.Name),
		Description: strings.TrimSpace(in.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.cats.Create(ctx, c)
}

type UpdateCategoryInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func (s *CategoryService) Update(ctx context.Context, id string, in UpdateCategoryInput) (models.Category, error) {
	patch := repositories.CategoryPatch{}
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		patch.Name = &n
	}
	if in.Description != nil {
		d := strings.TrimSpace(*in.Description)
		patch.Description = &d
	}
	now := s.clk.Now().Format(time.RFC3339)
	patch.UpdatedAt = &now
	return s.cats.Update(ctx, id, patch)
}

func (s *CategoryService) Delete(ctx context.Context, id string) error {
	return s.cats.Delete(ctx, id)
}
