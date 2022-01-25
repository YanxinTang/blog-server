package controller

import (
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/model"
	"github.com/YanxinTang/blog-server/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetStorages(c *gin.Context) {
	storages, err := model.GetStorages()
	if err != nil {
		log.Error("failed to get storages", zap.Error(err))
		c.Error(e.New(http.StatusBadRequest, "存储列表获取失败"))
		return
	}
	c.JSON(http.StatusOK, storages)
}

func GetStorage(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("categoryID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "ID 应该是整数"))
		return
	}
	storage, err := model.GetStorage(storageID)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "存储信息获取失败"))
	}
	c.JSON(http.StatusOK, storage)
}

type createStorageBody struct {
	Name      string `json:"name" binding:"required"`
	SecretID  string `json:"secretID" binding:"required"`
	SecretKey string `json:"secretKey" binding:"required"`
	Token     string `json:"token"`
	Region    string `json:"region" binding:"required"`
	Endpoint  string `json:"endpoint" binding:"required"`
	Bucket    string `json:"bucket" binding:"required"`
	Capacity  int64  `json:"capacity" binding:"required"`
}

func CreateStorage(c *gin.Context) {
	var body createStorageBody
	if err := c.BindJSON(&body); err != nil {
		return
	}
	var storage model.Storage
	storage.Name = body.Name
	storage.SecretID = body.SecretID
	storage.SecretKey = body.SecretKey
	storage.Token = body.Token
	storage.Region = body.Region
	storage.Endpoint = body.Endpoint
	storage.Bucket = body.Bucket

	storage, err := model.CreateStorage(storage)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "存储创建失败"))
		return
	}

	c.JSON(http.StatusOK, storage)
}

type updateStorageBody struct {
	Name      string `json:"name" binding:"required"`
	SecretID  string `json:"secretID" binding:"required"`
	SecretKey string `json:"secretKey" binding:"required"`
	Token     string `json:"token"`
	Region    string `json:"region" binding:"required"`
	Endpoint  string `json:"endpoint" binding:"required"`
	Bucket    string `json:"bucket" binding:"required"`
	Capacity  int64  `json:"capacity" binding:"required"`
}

func UpdateStorage(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("storageID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "ID 应该是整数"))
		return
	}

	var body updateStorageBody
	if err := c.BindJSON(&body); err != nil {
		return
	}

	var storage model.Storage
	storage.ID = storageID
	storage.Name = body.Name
	storage.SecretID = body.SecretID
	storage.SecretKey = body.SecretKey
	storage.Token = body.Token
	storage.Region = body.Region
	storage.Endpoint = body.Endpoint
	storage.Bucket = body.Bucket
	storage.Capacity = body.Capacity

	err = model.UpdateStorage(storage)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到此存储"))
		return
	}
	c.JSON(http.StatusOK, storage)
}

func DeleteStorage(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("storageID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "ID 应该是整数"))
		return
	}
	err = model.DeleteStorage(storageID)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到此存储"))
		return
	}
	c.Status(http.StatusOK)
}

func GetStorageObjects(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("storageID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到对存储库"))
		return
	}
	search := c.Param("search")
	nextContinuationToken := c.Param("nextContinuationToken")
	listObjectsV2Input := s3.ListObjectsV2Input{
		Prefix: aws.String(search),
	}
	// 使用 aws.String("") 初始化 ContinuationToken 会报 token 错误
	// 因此在这里判断不为空再设置 token
	if nextContinuationToken != "" {
		listObjectsV2Input.SetContinuationToken(nextContinuationToken)
	}
	listObjectsOutput, err := service.StorageListObjects(storageID, listObjectsV2Input)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, listObjectsOutput)
}

type DeleteStorageObjectQuery struct {
	Key string `form:"key" binding:"required"`
}

func PutStorageObject(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("storageID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到对存储库"))
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "获取上传文件失败"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "打开文件失败"))
		return
	}

	pubObjectInput := s3.PutObjectInput{
		Body: file,
		Key:  aws.String(fileHeader.Filename),
	}
	pubObjectInput.SetBody(file)
	pubObjectOutput, err := service.StoragePubObject(storageID, pubObjectInput)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, pubObjectOutput)
}

func DeleteStorageObject(c *gin.Context) {
	storageID, err := strconv.ParseUint(c.Param("storageID"), 10, 64)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到对存储库"))
		return
	}
	var deleteStorageObjectQuery DeleteStorageObjectQuery
	if err := c.BindQuery(&deleteStorageObjectQuery); err != nil {
		c.Error(err)
		return
	}

	deleteObjectInput := s3.DeleteObjectInput{
		Key: aws.String(deleteStorageObjectQuery.Key),
	}

	deleteObjectOutput, err := service.StorageDeleteObject(storageID, deleteObjectInput)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, deleteObjectOutput)
}
