package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	setting "file_storage_service/infrastructure/config"
)

func NewS3Client(common *setting.Config) *s3.Client {
	// ถ้ารันใน Docker Compose network เดียวกัน ใช้ "minio"
	// ถ้ารันนอก Docker ใช้ "localhost"
	endpoint := fmt.Sprintf("http://%s:%d", common.Storage.Host, common.Storage.Port)

	accessKey := common.Storage.AccessKey
	secretKey := common.Storage.SecretKey
	region := common.Storage.Region

	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // 🔑 สำคัญ
	})
}
