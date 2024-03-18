package services

import (
	"context"

	"github.com/Armunz/learn-docker/internal/entity"
	"github.com/Armunz/learn-docker/internal/model"
	"github.com/Armunz/learn-docker/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	CreateUser(ctx context.Context, request model.UserCreateRequest) error
	GetListUser(ctx context.Context, request model.UserListRequest) ([]model.UserResponse, int64, int64, error)
	GetUserDetail(ctx context.Context, userID string) (model.UserResponse, error)
	UpdateUser(ctx context.Context, userID string, request model.UserUpdateRequest) error
	DeleteUser(ctx context.Context, userID string) error
}

type serviceImpl struct {
	repository   repositories.Repository
	defaultLimit int
}

func NewService(repository repositories.Repository, defaultLimit int) Service {
	return &serviceImpl{
		repository:   repository,
		defaultLimit: defaultLimit,
	}
}

// CreateUser implements Service.
func (s *serviceImpl) CreateUser(ctx context.Context, request model.UserCreateRequest) error {
	user := entity.User{
		UserID: uuid.NewString(),
		Name:   request.Name,
		Age:    request.Age,
	}

	return s.repository.Create(ctx, user)
}

// DeleteUser implements Service.
func (s *serviceImpl) DeleteUser(ctx context.Context, userID string) error {
	return s.repository.Delete(ctx, userID)
}

// GetListUser implements Service.
func (s *serviceImpl) GetListUser(ctx context.Context, request model.UserListRequest) ([]model.UserResponse, int64, int64, error) {
	limit := request.Limit
	if limit == 0 {
		limit = s.defaultLimit
	}

	var offset int
	if request.Page > 0 {
		offset = (request.Page - 1) * limit
	}

	var count int64
	var users []entity.User

	g, groupCtx := errgroup.WithContext(ctx)
	// get count data
	g.Go(func() (err error) {
		count, err = s.repository.Count(groupCtx)
		return err
	})

	// get users
	g.Go(func() (err error) {
		users, err = s.repository.Get(groupCtx, limit, offset)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, 0, err
	}

	// count total pages
	var totalPages int64
	if limit > 0 {
		totalPages = count / int64(limit)
		if count%int64(limit) != 0 {
			totalPages++
		}
	}

	response := make([]model.UserResponse, len(users))
	for i, u := range users {
		response[i] = model.UserResponse{
			Name: u.Name,
			Age:  u.Age,
		}
	}

	return response, count, totalPages, nil
}

// GetUserDetail implements Service.
func (s *serviceImpl) GetUserDetail(ctx context.Context, userID string) (model.UserResponse, error) {
	user, err := s.repository.GetById(ctx, userID)
	if err != nil {
		return model.UserResponse{}, err
	}

	response := model.UserResponse{
		Name: user.Name,
		Age:  user.Age,
	}

	return response, nil
}

// UpdateUser implements Service.
func (s *serviceImpl) UpdateUser(ctx context.Context, userID string, request model.UserUpdateRequest) error {
	user, err := s.repository.GetById(ctx, userID)
	if err != nil {
		return err
	}

	user.Name = request.Name
	user.Age = request.Age

	return s.repository.Update(ctx, user)
}
