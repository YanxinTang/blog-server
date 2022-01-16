package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YanxinTang/blog-server/e"
	"github.com/YanxinTang/blog-server/model"
	"github.com/gin-gonic/gin"
)

func GetStorages(c *gin.Context) {
	storages, err := model.GetStorages()
	if err != nil {
		log.Println(err)
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
	Token     string `json:"token" binding:"required"`
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
		log.Println(err)
		c.Error(e.New(http.StatusBadRequest, "存储创建失败"))
		return
	}

	c.JSON(http.StatusOK, storage)
}

type updateStorageBody struct {
	Name      string `json:"name" binding:"required"`
	SecretID  string `json:"secretID" binding:"required"`
	SecretKey string `json:"secretKey" binding:"required"`
	Token     string `json:"token" binding:"required"`
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
		log.Println(err)
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
		log.Println(err)
		c.Error(e.New(http.StatusNotFound, "找不到此存储"))
		return
	}
	c.Status(http.StatusOK)
}
