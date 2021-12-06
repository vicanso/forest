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
	"context"

	"github.com/vicanso/forest/ent"
	"github.com/vicanso/forest/ent/file"
	"github.com/vicanso/forest/helper"
	"github.com/vicanso/forest/validate"
)

type entStorage struct {
	client *ent.Client
}

func mustNewEntStorage() FileStorage {
	return &entStorage{
		client: helper.EntGetClient(),
	}
}

func convertToFile(data *ent.File) *File {
	return &File{
		Bucket:      data.Bucket,
		Filename:    data.Filename,
		ContentType: data.ContentType,
		Size:        data.Size,
		Metadata:    *data.Metadata,
		Creator:     data.Creator,
		Data:        data.Data,
	}
}

// Get gets file from ent(mysql or postgres)
func (e *entStorage) Get(ctx context.Context, bucket, filename string) (*File, error) {
	result, err := e.client.File.Query().
		Where(file.Bucket(bucket)).
		Where(file.Filename(filename)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return convertToFile(result), nil
}

// Put puts file to ent(mysql or postgres)
func (e *entStorage) Put(ctx context.Context, data File) error {
	err := validate.Struct(&data)
	if err != nil {
		return err
	}
	// 如果指定了id，则是更新
	if data.ID != 0 {
		_, err = e.client.File.UpdateOneID(data.ID).
			SetContentType(data.ContentType).
			SetSize(data.Size).
			SetMetadata(&data.Metadata).
			SetData(data.Data).
			Save(ctx)
		if err != nil {
			return err
		}
	}
	_, err = e.client.File.Create().
		SetBucket(data.Bucket).
		SetFilename(data.Filename).
		SetContentType(data.ContentType).
		SetSize(data.Size).
		SetMetadata(&data.Metadata).
		SetCreator(data.Creator).
		SetData(data.Data).
		Save(ctx)
	return err
}

// Query gets the files from ent(mysql or postgres)
func (e *entStorage) Query(ctx context.Context, param FileFilterParams) ([]*File, error) {
	return nil, nil
}

// Count counts the files from ent(mysql or postgres)
func (e *entStorage) Count(ctx context.Context, params FileFilterParams) (int64, error) {
	return -1, nil
}
