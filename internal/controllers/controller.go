package controllers

import (
	"context"
	"time"

	"github.com/Armunz/learn-docker/internal/model"
	"github.com/Armunz/learn-docker/internal/services"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type resource struct {
	service  services.Service
	validate *validator.Validate
	timeout  int
}

func RegisterHandlers(r fiber.Router, service services.Service, validate *validator.Validate, timeout int) {
	res := resource{
		service:  service,
		validate: validate,
		timeout:  timeout,
	}

	r.Get("/hello", res.Hello)
	r.Post("/user", res.CreateUser)
	r.Get("/user", res.GetListUser)
	r.Get("/user/:id", res.GetUserDetail)
	r.Put("/user/:id", res.UpdateUser)
	r.Delete("/user/:id", res.DeleteUser)
}

func (r *resource) Hello(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}

func (r *resource) CreateUser(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	var request model.UserCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	if err := r.validate.Struct(request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	if err := r.service.CreateUser(c.UserContext(), request); err != nil {
		return model.Response(c, fiber.StatusInternalServerError, err)
	}

	return model.Response(c, fiber.StatusCreated, nil)
}

func (r *resource) GetListUser(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	var request model.UserListRequest
	if err := c.QueryParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	response, totalData, totalPage, err := r.service.GetListUser(c.UserContext(), request)
	if err != nil {
		return model.Response(c, fiber.StatusInternalServerError, nil)
	}

	responsePage := model.ResponsePage{
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return model.Response(c, fiber.StatusOK, response, responsePage)
}

func (r *resource) GetUserDetail(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	userID := c.Params("id")
	if userID == "" {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	response, err := r.service.GetUserDetail(c.UserContext(), userID)
	if err != nil {
		return model.Response(c, fiber.StatusInternalServerError, nil)
	}

	return model.Response(c, fiber.StatusOK, response)
}

func (r *resource) UpdateUser(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	var request model.UserUpdateRequest
	if err := c.BodyParser(&request); err != nil {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	if err := r.validate.Struct(request); err != nil {
		return model.Response(c, fiber.StatusBadRequest)
	}

	userID := c.Params("id")
	if userID == "" {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	if err := r.service.UpdateUser(c.UserContext(), userID, request); err != nil {
		return model.Response(c, fiber.StatusInternalServerError, nil)
	}

	return model.Response(c, fiber.StatusOK, nil)
}

func (r *resource) DeleteUser(c *fiber.Ctx) error {
	// set timeout
	timeout, cancel := context.WithTimeout(c.UserContext(), time.Duration(r.timeout)*time.Second)
	defer cancel()
	c.SetUserContext(timeout)

	userID := c.Params("id")
	if userID == "" {
		return model.Response(c, fiber.StatusBadRequest, nil)
	}

	if err := r.service.DeleteUser(c.UserContext(), userID); err != nil {
		return model.Response(c, fiber.StatusInternalServerError, nil)
	}

	return model.Response(c, fiber.StatusOK, nil)
}
