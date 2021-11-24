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
	"net/http"
)

const creatorField = "creator"

// 文件
type File struct {
	// 文件id
	ID     int    `json:"id"`
	Bucket string `json:"bucket" validate:"required"`
	// 文件名
	Filename string `json:"filename" validate:"required"`
	// 类型
	ContentType string `json:"contentType" validate:"required"`
	// 大小
	Size int64 `json:"size" validate:"required"`
	// metadata
	Metadata http.Header `json:"metadata"`
	// 创建者
	Creator string `json:"creator" validate:"required"`
	// 数据
	Data []byte `json:"data" validate:"required"`
}
type FileFilterParams struct {
	// 筛选的字段
	Fields string `json:"fields"`
	// 数量
	Limit int `json:"limit"`
	// 偏移量
	Offset int `json:"offset"`
}

type FileStorage interface {
	Get(ctx context.Context, bucket, filename string) (*File, error)
	Put(ctx context.Context, file File) error
	Query(ctx context.Context, params FileFilterParams) ([]*File, error)
	Count(ctx context.Context, params FileFilterParams) (int64, error)
}

var minioStorageClient = mustNewMinioStorage()
var entStorageClient = mustNewEntStorage()

func Minio() FileStorage {
	return minioStorageClient
}

func Ent() FileStorage {
	return entStorageClient
}
