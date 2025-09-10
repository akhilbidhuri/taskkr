package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/akhilbidhuri/taskkr/internal/model"
	"github.com/akhilbidhuri/taskkr/internal/service"
	"github.com/akhilbidhuri/taskkr/internal/utils"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetTask)
	r.Post("/", h.CreateTask)
	r.Get("/", h.ListTasks)
	r.Put("/{id}", h.UpdateTask)
	r.Delete("/{id}", h.DeleteTask)
	return r
}

// GetTask godoc
// @Summary Get single task based on id query param
// @Description Get single task based on id if present
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int false "ID filter"
// @Success 200 {object} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "", err)
		return
	}
	if task == nil {
		utils.Error(w, http.StatusNotFound, "Task not found", nil)
		return
	}
	utils.Success(w, http.StatusOK, "", task)
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a task with title, description, etc.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task body model.Task true "Task info"
// @Success 201 {object} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	if task.UserID == 0 || task.Title == "" {
		utils.Error(w, http.StatusBadRequest, "missing values in body", nil)
		return
	}
	err := h.service.Create(r.Context(), &task)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "", err)
		return
	}
	utils.Success(w, http.StatusCreated, "", task)
}

// GetTasks godoc
// @Summary Get list of tasks
// @Description Get all tasks with pagination and optional filtering
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param status query string false "Filter by status" Enums(pending, in_process, completed)
// @Param title query string false "Title filter"
// @Param page query string false "Page filter"
// @Param page_size query string false "PageSize filter"
// @Success 200 {array} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	filter := &model.TaskFilter{
		Page:     1,
		PageSize: 10,
	}

	if params.Get("status") != "" {
		status, err := getStatus(params.Get("status"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "", err)
			return
		}
		filter.Status = status
	}
	if params.Get("title") != "" {
		filter.Title = params.Get("title")
	}
	if params.Get("page") != "" {
		page, err := strconv.ParseUint(params.Get("page"), 10, 32)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid page value", nil)
			return
		}
		filter.Page = uint(page)
	}
	if params.Get("page_size") != "" {
		pageSize, err := strconv.ParseUint(params.Get("page_size"), 10, 32)
		if err != nil || pageSize > 100 {
			utils.Error(w, http.StatusBadRequest, "Invalid page value", nil)
			return
		}
		filter.PageSize = uint(pageSize)
	}
	tasks, total, err := h.service.List(r.Context(), filter)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "", err)
		return
	}

	resp := map[string]interface{}{
		"total": total,
		"tasks": tasks,
	}
	utils.Success(w, http.StatusOK, "", resp)
}

// UpdateTasks godoc
// @Summary Update a task
// @Description UPdate a task with given ID and values
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path int true "ID filter"
// @Param task body model.UpdateTask true "Task update info"
// @Success 202 {array} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updateTask model.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&updateTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if updateTask.Status != "" {
		status, err := getStatus(string(updateTask.Status))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "", err)
			return
		}
		updateTask.Status = status
	}
	task, err := h.service.Update(r.Context(), id, &updateTask)
	if err != nil {
		if err == utils.NoEntryError {
			utils.Error(w, http.StatusNotFound, "", err)
			return
		}
		utils.Error(w, http.StatusInternalServerError, "", err)
		return
	}
	utils.Success(w, http.StatusAccepted, "", task)
}

// DeleteTasks godoc
// @Summary Delete a task
// @Description Delete a task with given ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string false "ID filter"
// @Success 202 {array} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if err == utils.NoEntryError {
			utils.Error(w, http.StatusNotFound, "", err)
			return
		}
		utils.Error(w, http.StatusInternalServerError, "", err)
		return
	}
	utils.Success(w, http.StatusNoContent, "", nil)
}

func getStatus(statusStr string) (model.TaskStatus, error) {
	switch model.TaskStatus(statusStr) {
	case model.StatusPending, model.StatusCompleted, model.StatusInProcess:
		return model.TaskStatus(statusStr), nil
	}
	return "", errors.New("Invlaid status value")
}
