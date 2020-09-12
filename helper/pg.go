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

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/ent"
)

var (
	client *ent.Client
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
	// defer db.Close()

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

// GetEntClient get ent client
func GetEntClient() *ent.Client {
	return client
}

// InitSchemaAndHook 初始化schema与hook函数
func InitSchemaAndHook() (err error) {
	err = client.Schema.Create(context.Background())
	if err != nil {
		return
	}
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			// start := time.Now()
			// defer func() {
			// 	log.Printf("Op=%s\tType=%s\tTime=%s\tConcreteType=%T\n", m.Op(), m.Type(), time.Since(start), m)
			// }()
			return next.Mutate(ctx, m)
		})
	})
	return
}
