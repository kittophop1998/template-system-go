package model

type Attachment struct {
	ID       string `gorm:"primaryKey" json:"id"`
	FileName string `gorm:"column:file_name" json:"file_name"`
	Size     int64  `gorm:"column:size" json:"size"`
}
