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

package middleware

import (
	"github.com/vicanso/elton"
	session "github.com/vicanso/elton-session"
	"github.com/vicanso/forest/cache"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/util"
)

// NewSession new session middleware
func NewSession() elton.Handler {
	store := &cache.RedisSessionStore{
		Prefix: "ss-",
	}
	scf := config.GetSessionConfig()
	return session.NewByCookie(session.CookieConfig{
		Store:   store,
		Signed:  true,
		Expired: scf.TTL,
		GenID: func() string {
			return util.GenUlid()
		},
		Name:     scf.Key,
		Path:     scf.CookiePath,
		MaxAge:   int(scf.TTL.Seconds()),
		HttpOnly: true,
	})
}
