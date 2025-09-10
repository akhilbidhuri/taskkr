package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akhilbidhuri/taskkr/internal/model"
	"github.com/akhilbidhuri/taskkr/internal/service"
	"github.com/akhilbidhuri/taskkr/utils"

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
	//r.Get("/", h.ListTasks)
	// r.Put("/{id}", h.UpdateTask)
	// r.Delete("/{id}", h.DeleteTask)
	return r
}

// GetTask godoc
// @Summary Get single task based on id query param
// @Description Get single task based on id if present
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id query string false "ID filter"
// @Success 200 {object} model.Task
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/tasks/{id} [get]
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
	utils.Success(w, http.StatusFound, "", task)
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
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	err := h.service.Create(r.Context(), &task)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "", err)
		return
	}
	utils.Success(w, http.StatusCreated, "", task)
}
