package controller

import (
	"io"
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/internal/app/common"
	"github.com/YanxinTang/blog-server/internal/app/service"
	"github.com/YanxinTang/blog-server/internal/pkg/e"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/YanxinTang/blog-server/internal/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// cover writer to writerat
type FakerWriteAt struct {
	w io.Writer
}

func (w *FakerWriteAt) WriteAt(p []byte, off int64) (n int, err error) {
	return w.w.Write(p)
}

type CreateStorageReqBody struct {
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
	var body CreateStorageReqBody
	if err := c.BindJSON(&body); err != nil {
		return
	}

	csi := model.CreateStorageInput{
		Name:      body.Name,
		SecretID:  body.SecretID,
		SecretKey: body.SecretKey,
		Token:     body.Token,
		Region:    body.Region,
		Endpoint:  body.Endpoint,
		Bucket:    body.Bucket,
	}

	storage, err := model.CreateStorage(common.Context, common.Client)(csi)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "存储创建失败"))
		return
	}

	c.JSON(http.StatusOK, storage)
}

type UpdateStorageReqBody struct {
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
	storageID, err := strconv.Atoi(c.Param("storageID"))
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "ID 应该是整数"))
		return
	}

	var body UpdateStorageReqBody
	if err := c.BindJSON(&body); err != nil {
		return
	}

	var usi model.UpdateStorageInput
	usi.ID = storageID
	usi.Name = body.Name
	usi.SecretID = body.SecretID
	usi.SecretKey = body.SecretKey
	usi.Token = body.Token
	usi.Region = body.Region
	usi.Endpoint = body.Endpoint
	usi.Bucket = body.Bucket
	usi.Capacity = body.Capacity

	s, err := model.UpdateStorage(common.Context, common.Client)(usi)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到此存储"))
		return
	}
	c.JSON(http.StatusOK, s)
}

func DeleteStorage(c *gin.Context) {
	storageID, err := strconv.Atoi(c.Param("storageID"))
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	err = model.DeleteStorage(common.Context, common.Client)(storageID)
	if err != nil {
		c.Error(e.New(http.StatusNotFound, "找不到此存储"))
		return
	}
	c.Status(http.StatusOK)
}

func GetStorageObjects(c *gin.Context) {
	storageID, err := strconv.Atoi(c.Param("storageID"))
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
	storageID, err := strconv.Atoi(c.Param("storageID"))
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
	storageID, err := strconv.Atoi(c.Param("storageID"))
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

func GetStorages(c *gin.Context) {
	storages, err := model.GetStorages(common.Context, common.Client)()
	if err != nil {
		log.Error("failed to get storages", zap.Error(err))
		c.Error(e.New(http.StatusBadRequest, "存储列表获取失败"))
		return
	}
	c.JSON(http.StatusOK, storages)
}

func GetStorage(c *gin.Context) {
	storageID, err := strconv.Atoi(c.Param("storageID"))
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "ID 应该是整数"))
		return
	}
	storage, err := model.GetStorage(common.Context, common.Client)(storageID)
	if err != nil {
		c.Error(e.New(http.StatusBadRequest, "存储信息获取失败"))
	}
	c.JSON(http.StatusOK, storage)
}

func GetStorageObject(c *gin.Context) {
	storageID, err := strconv.Atoi(c.Param("storageID"))
	if err != nil {
		c.Error(e.ERROR_BAD_REQUEST)
		return
	}
	d, err := service.GetStorageDownloader(storageID)
	if err != nil {
		c.Error(err)
		return
	}
	key := c.Param("key")
	getObjectInput := s3.GetObjectInput{
		Bucket: aws.String(d.Storage.Bucket),
		Key:    aws.String(key),
	}
	writerAt := FakerWriteAt{c.Writer}
	defer (func() {
		_, err := d.Downloader.Download(&writerAt, &getObjectInput)
		if err != nil {
			log.Warn(
				"failed to download",
				zap.Int("storageID", storageID),
				zap.String("key", key),
				zap.Error(err),
			)
		}
	})()
	c.Status(http.StatusOK)
}
