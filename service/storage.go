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

func GetStorageService(storageID uint64) (*StorageService, e.ApiError) {
	conf, err := model.GetStorage(storageID)
	if err != nil {
		log.Println("model.GetStorage: ", err)
		return nil, e.New(http.StatusBadRequest, "获取存储失败")
	}
	sess, err := newS3Session(conf)
	if err != nil {
		log.Println("newS3Session: ", err)
		return nil, e.New(http.StatusInternalServerError, "初始化 S3 会话失败")
	}
	svc := StorageService{
		Storage: &conf,
		S3:      s3.New(sess),
	}
	return &svc, nil
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

func StorageListObjects(storageID uint64, listObjectsV2Input s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, e.ApiError) {
	conf, err := model.GetStorage(storageID)
	if err != nil {
		log.Println("model.GetStorage: ", err)
		return nil, e.New(http.StatusBadRequest, "获取存储失败")
	}
	sess, err := newS3Session(conf)
	if err != nil {
		log.Println("newS3Session: ", err)
		return nil, e.New(http.StatusInternalServerError, "初始化 S3 会话失败")
	}
	s3svc := s3.New(sess)
	listObjectsV2Input.SetBucket(conf.Bucket)
	listObjectsV2Input.SetMaxKeys(100)

	listObjectsV2Output, err := s3svc.ListObjectsV2(&listObjectsV2Input)
	if err != nil {
		log.Println("s3svc.ListObjectsV2: ", err)
		return nil, e.New(http.StatusBadRequest, "获取存储文件列表失败")
	}
	return listObjectsV2Output, nil
}

func StoragePubObject(storageID uint64, putObjectInput s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	svc, apiError := GetStorageService(storageID)
	if apiError != nil {
		return nil, apiError
	}
	putObjectInput.SetBucket(svc.Storage.Bucket)
	pubObjectOutput, err := svc.S3.PutObject(&putObjectInput)
	if err != nil {
		log.Println(err)
		return nil, e.New(http.StatusBadRequest, "上传文件失败")
	}
	return pubObjectOutput, nil
}

func StorageDeleteObject(storageID uint64, deleteObjectInput s3.DeleteObjectInput) (*s3.DeleteObjectOutput, e.ApiError) {
	svc, apiError := GetStorageService(storageID)
	if apiError != nil {
		return nil, apiError
	}
	deleteObjectInput.SetBucket(svc.Storage.Bucket)
	deleteObjectOutput, err := svc.S3.DeleteObject(&deleteObjectInput)
	if err != nil {
		e.New(http.StatusBadRequest, "删除文件失败")
	}
	return deleteObjectOutput, nil
}
