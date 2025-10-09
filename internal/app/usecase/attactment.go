package usecase

import (
	"bytes"
	"context"
	"file_storage_service/internal/domain/model"
	"file_storage_service/internal/domain/repository"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AttachmentUseCase struct {
	S3     *s3.Client
	Bucket string
	Repo   repository.AttachmentRepository
}

func NewAttachmentUseCase(s3Client *s3.Client, bucket string, repo repository.AttachmentRepository) *AttachmentUseCase {
	return &AttachmentUseCase{
		S3:     s3Client,
		Bucket: bucket,
		Repo:   repo,
	}
}

func (fs *AttachmentUseCase) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	key := header.Filename

	// Upload file (streaming, no full buffer)
	_, err := fs.S3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &fs.Bucket,
		Key:         &key,
		Body:        file,
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Retrieve file metadata
	headResp, err := fs.S3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get file metadata: %w", err)
	}

	// ===== Save attachment metadata =====
	attachmentMeta := model.Attachment{
		FileName: header.Filename,
		Size:     aws.ToInt64(headResp.ContentLength),
	}

	if err := fs.Repo.SaveAttachName(&attachmentMeta); err != nil {
		return "", fmt.Errorf("failed to save attachment metadata: %w", err)
	}
	// =====================================

	return key, nil
}

func (fs *AttachmentUseCase) DownloadFile(key string) ([]byte, error) {
	resp, err := fs.S3.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (fs *AttachmentUseCase) DeleteFile(key string) error {
	_, err := fs.S3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &fs.Bucket,
		Key:    &key,
	})
	return err
}
