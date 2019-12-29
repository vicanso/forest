// Copyright 2019 tree xie
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

package helper

import (
	"regexp"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
	"go.uber.org/zap"
)

var (
	pgClient *gorm.DB
	pgSlow   = time.Second
)

func pgAddStartedAt(scope *gorm.Scope) {
	scope.InstanceSet(startedAtKey, time.Now())
}
func pgStats(category string) func(*gorm.Scope) {
	return func(scope *gorm.Scope) {
		value, ok := scope.InstanceGet(startedAtKey)
		if !ok {
			return
		}
		startedAt, ok := value.(time.Time)
		if !ok {
			return
		}
		use := time.Since(startedAt)
		db := scope.DB()
		if time.Since(startedAt) > pgSlow || db.Error != nil {
			message := ""
			if db.Error != nil {
				message = db.Error.Error()
			}
			logger.Info("pg process slow or error",
				zap.String("table", scope.TableName()),
				zap.String("category", category),
				zap.String("use", use.String()),
				zap.Int64("rowsAffected", db.RowsAffected),
				zap.String("error", message),
			)
			tags := map[string]string{
				"table":    scope.TableName(),
				"category": category,
			}
			fields := map[string]interface{}{
				"use":          use.Milliseconds(),
				"rowsAffected": db.RowsAffected,
				"error":        message,
			}
			GetInfluxSrv().Write(cs.MeasurementPG, fields, tags)
		}
	}
}

func init() {
	str := config.GetPostgresConnectString()
	reg := regexp.MustCompile(`password=\S*`)
	maskStr := reg.ReplaceAllString(str, "password=***")
	logger.Info("connect to pg",
		zap.String("args", maskStr),
	)
	db, err := gorm.Open("postgres", str)
	if err != nil {
		panic(err)
	}
	db.SetLogger(log.PGLogger())
	db.Callback().Query().Before("gorm:query").Register("stats:beforeQuery", pgAddStartedAt)
	db.Callback().Query().After("gorm:query").Register("stats:afterQuery", pgStats("query"))
	db.Callback().Update().Before("gorm:update").Register("stats:beforeUpdate", pgAddStartedAt)
	db.Callback().Update().After("gorm:update").Register("stats:afterUpdate", pgStats("update"))

	pgClient = db
}

// PGCreate pg create
func PGCreate(data interface{}) (err error) {
	err = pgClient.Create(data).Error
	return
}

// PGGetClient pg client
func PGGetClient() *gorm.DB {
	return pgClient
}

// PGFormatOrder format order
func PGFormatOrder(sort string) string {
	arr := strings.Split(sort, ",")
	newSort := []string{}
	for _, item := range arr {
		if item[0] == '-' {
			newSort = append(newSort, strcase.ToSnake(item[1:])+" desc")
		} else {
			newSort = append(newSort, strcase.ToSnake(item))
		}
	}
	return strings.Join(newSort, ",")
}

// PGFormatSelect format select
func PGFormatSelect(fields string) string {
	return strcase.ToSnake(fields)
}
