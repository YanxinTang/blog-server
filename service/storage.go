package service

import (
	"log"
	"net/http"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type StorageService struct {
	Storage *model.Storage
	S3      *s3.S3
}

func newS3Session(conf model.Storage) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(conf.SecretID, conf.SecretKey, conf.Token)
	config := &aws.Config{
		Region:           &conf.Region,
		Endpoint:         &conf.Endpoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	}
	return session.NewSession(config)
}

func GetStorageServices() ([]StorageService, e.ApiError) {
	var services []StorageService
	storages, err := model.GetStorages()
	if err != nil {
		log.Println("存储列表获取失败: ", err)
		return nil, e.New(http.StatusBadRequest, "存储列表获取失败")
	}
	for _, conf := range storages {
		sess, err := newS3Session(conf)
		if err != nil {
			log.Println("newS3Session: ", err)
			return nil, e.New(http.StatusInternalServerError, "初始化 S3 会话失败")
		}
		svc := StorageService{
			Storage: &conf,
			S3:      s3.New(sess),
		}
		services = append(services, svc)
	}
	return services, nil
}

func GetStorageUsage(svc StorageService) (int64, e.ApiError) {
	listBucketsInput := s3.ListBucketsInput{}
	listBucketsOutput, err := svc.S3.ListBuckets(&listBucketsInput)
	if err != nil {
		log.Println("获取存储桶列表失败: ", err)
		return 0, e.New(http.StatusBadRequest, "获取存储桶列表失败")
	}
	var usage int64
	for _, bucket := range listBucketsOutput.Buckets {
		bucketUsage, err := GetStorageBucketUsage(svc, *bucket.Name)
		if err != nil {
			return 0, err
		}
		usage += bucketUsage
	}
	return usage, nil
}

func GetStorageBucketUsage(svc StorageService, name string) (int64, e.ApiError) {
	listObjectsV2Input := s3.ListObjectsV2Input{
		Bucket: aws.String(name),
	}
	var usage int64
	if err := svc.S3.ListObjectsV2Pages(
		&listObjectsV2Input,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range page.Contents {
				usage += *obj.Size
			}
			return !lastPage
		},
	); err != nil {
		return 0, e.New(http.StatusBadRequest, "获取存储列表失败")
	}
	return usage, nil
}
