// Copyright 2020 tree xie
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
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/ent"
)

var (
	client *ent.Client
	driver *entsql.Driver
)

func init() {
	postgresConfig := config.GetPostgresConfig()
	c, err := open(postgresConfig.URI)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err := c.Schema.Create(ctx); err != nil {
		panic(err)
	}
	client = c
}

// open new connection
func open(databaseUrl string) (*ent.Client, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}

	// Create an ent.Driver from `db`.
	driver = entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(driver)), nil
}

// GetEntClient get ent client
func GetEntClient() *ent.Client {
	return client
}

// EntPing ent driver ping
func EntPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return driver.DB().PingContext(ctx)
}

// InitSchemaAndHook 初始化schema与hook函数
func InitSchemaAndHook() (err error) {
	err = client.Schema.Create(context.Background())
	if err != nil {
		return
	}
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			start := time.Now()
			defer func() {
				for _, name := range m.Fields() {
					fmt.Println(m.Field(name))
				}
				log.Printf("Op=%s\tType=%s\tTime=%s\tConcreteType=%T\n", m.Op(), m.Type(), time.Since(start), m)
			}()
			return next.Mutate(ctx, m)
		})
	})
	return
}
