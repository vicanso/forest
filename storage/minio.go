// Copyright 2021 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/validate"
)

func mustNewMinioStorage() FileStorage {
	minioConfig := config.MustGetMinioConfig()
	c, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Secure: minioConfig.SSL,
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
	})
	if err != nil {
		panic(err)
	}
	return &minioStorage{
		client: c,
	}
}

type minioStorage struct {
	client *minio.Client
}

// Get gets file from minio
func (m *minioStorage) Get(ctx context.Context, bucket, filename string) (*File, error) {
	obj, err := m.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	statsInfo, err := obj.Stat()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return &File{
		Bucket:      bucket,
		Filename:    filename,
		ContentType: statsInfo.ContentType,
		Size:        statsInfo.Size,
		Metadata:    statsInfo.Metadata,
		// Creator:     statsInfo.Owner.DisplayName,
		Data: data,
	}, nil
}

// Put puts file to minio
func (m *minioStorage) Put(ctx context.Context, file File) error {
	err := validate.Struct(&file)
	if err != nil {
		return err
	}
	r := bytes.NewReader(file.Data)
	size := int64(len(file.Data))
	metadata := make(map[string]string)
	for key, values := range file.Metadata {
		metadata[key] = strings.Join(values, ",")
	}

	_, err = m.client.PutObject(
		ctx,
		file.Bucket,
		file.Filename,
		r,
		size,
		minio.PutObjectOptions{
			ContentType:  file.ContentType,
			UserMetadata: metadata,
		},
	)
	return err
}

// Query gets the files from minio
func (m *minioStorage) Query(ctx context.Context, param FileFilterParams) ([]*File, error) {
	return nil, nil
}

// Count counts the files from minio
func (m *minioStorage) Count(ctx context.Context, params FileFilterParams) (int64, error) {
	return -1, nil
}
