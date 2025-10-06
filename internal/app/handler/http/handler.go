package http

import "file_storage_service/internal/app/usecase"

// ===== Handlers groups all handlers. =====
type Handlers struct {
	User       *UserHandler
	Attachment *AttachmentHandler
}

// ===== HandlerDependency groups all dependencies for handlers. =====
type HandlerDependency struct {
	UserUC       *usecase.UserUseCase
	AttachmentUC *usecase.AttachmentUseCase
}

// ===== InitializeHandlers initializes all handlers with their dependencies. =====
func InitializeHandlers(deps *HandlerDependency) *Handlers {
	return &Handlers{
		User:       &UserHandler{UserUC: deps.UserUC},
		Attachment: &AttachmentHandler{AttachmentUC: deps.AttachmentUC},
	}
}
