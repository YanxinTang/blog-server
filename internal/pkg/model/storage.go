package model

import (
	"context"

	"github.com/YanxinTang/blog-server/ent"
	"github.com/YanxinTang/blog-server/internal/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func GetStorages(ctx context.Context, client *ent.Client) func() ([]*ent.Storage, error) {
	return func() ([]*ent.Storage, error) {
		storages, err := client.Storage.Query().All(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failing getting all storages")
		}
		return storages, nil
	}
}

func GetStorage(ctx context.Context, client *ent.Client) func(id int) (*ent.Storage, error) {
	return func(id int) (*ent.Storage, error) {
		s, err := client.Storage.Get(ctx, id)
		if err != nil {
			return nil, errors.Wrapf(err, "failing getting storage[%d]", id)
		}
		return s, nil
	}
}

type CreateStorageInput struct {
	Name      string
	SecretID  string
	SecretKey string
	Token     string
	Region    string
	Endpoint  string
	Bucket    string
	Usage     int64
	Capacity  int64
}

func CreateStorage(ctx context.Context, client *ent.Client) func(CreateStorageInput) (*ent.Storage, error) {
	return func(csi CreateStorageInput) (*ent.Storage, error) {
		s, err := client.Storage.
			Create().
			SetName(csi.Name).
			SetSecretID(csi.SecretID).
			SetSecretKey(csi.SecretKey).
			SetToken(csi.Token).
			SetRegion(csi.Region).
			SetEndpoint(csi.Endpoint).
			SetBucket(csi.Bucket).
			SetCapacity(csi.Capacity).
			Save(ctx)
		if err != nil {
			log.Warn("failing creating storage", zap.Error(err))
			return nil, err
		}
		return s, nil
	}
}

type UpdateStorageInput struct {
	ID        int
	Name      string
	SecretID  string
	SecretKey string
	Token     string
	Region    string
	Endpoint  string
	Bucket    string
	Usage     int64
	Capacity  int64
}

func UpdateStorage(ctx context.Context, client *ent.Client) func(UpdateStorageInput) (*ent.Storage, error) {
	return func(usi UpdateStorageInput) (*ent.Storage, error) {
		updateStorageClient := client.Storage.UpdateOneID(usi.ID)
		if usi.Name != "" {
			updateStorageClient = updateStorageClient.SetName(usi.Name)
		}
		if usi.SecretID != "" {
			updateStorageClient = updateStorageClient.SetSecretID(usi.SecretID)
		}
		if usi.SecretKey != "" {
			updateStorageClient = updateStorageClient.SetSecretKey(usi.SecretKey)
		}
		if usi.Token != "" {
			updateStorageClient = updateStorageClient.SetToken(usi.Token)
		}
		if usi.Region != "" {
			updateStorageClient = updateStorageClient.SetRegion(usi.Region)
		}
		if usi.Endpoint != "" {
			updateStorageClient = updateStorageClient.SetEndpoint(usi.Endpoint)
		}
		if usi.Bucket != "" {
			updateStorageClient = updateStorageClient.SetBucket(usi.Bucket)
		}
		if usi.Usage != 0 {
			updateStorageClient = updateStorageClient.SetUsage(usi.Usage)
		}
		if usi.Capacity != 0 {
			updateStorageClient = updateStorageClient.SetCapacity(usi.Capacity)
		}
		s, err := updateStorageClient.Save(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "failing updating storage[%d]", usi.ID)
		}
		return s, nil
	}
}

func DeleteStorage(ctx context.Context, client *ent.Client) func(id int) error {
	return func(id int) error {
		err := client.Storage.DeleteOneID(id).Exec(ctx)
		if err != nil {
			return errors.Wrapf(err, "failing deleting storage[%d]", id)
		}
		return nil
	}
}
