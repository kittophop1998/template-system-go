package repository

import "file_storage_service/internal/domain/model"

type AttachmentRepository interface {
	SaveAttachName(attachment *model.Attachment) error
}
