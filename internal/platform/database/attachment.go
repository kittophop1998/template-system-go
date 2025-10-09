package database

import (
	"crypto/rand"
	"file_storage_service/internal/domain/model"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type AttachmentPostGres struct {
	db *gorm.DB
}

func NewAttachmentPostGres(db *gorm.DB) *AttachmentPostGres {
	return &AttachmentPostGres{db: db}
}

func (r *AttachmentPostGres) SaveAttachName(attachment *model.Attachment) error {
	newAttachment := &model.Attachment{
		ID:       ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String(),
		FileName: attachment.FileName,
		Size:     attachment.Size,
	}

	return r.db.Create(newAttachment).Error
}
