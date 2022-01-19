// Copyright 2022 tree xie
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

package service

import (
	"context"
	"runtime"
	"time"

	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
)

// https://eng.uber.com/how-we-saved-70k-cores-across-30-mission-critical-services/

type finalizer struct {
	ch  chan time.Time
	ref *finalizerRef
}
type finalizerRef struct {
	parent *finalizer
}

func finalizerHandler(f *finalizerRef) {
	select {
	case f.parent.ch <- time.Time{}:
	default:
		GetInfluxSrv().Write(cs.MeasurementEvent, map[string]string{
			cs.TagCategory: "gc",
		}, map[string]interface{}{
			cs.FieldCount: 1,
		})
		log.Info(context.Background()).
			Msg("gc")
	}
	runtime.SetFinalizer(f, finalizerHandler)
}

func NewTicker() *finalizer {
	f := &finalizer{
		ch: make(chan time.Time, 1),
	}
	f.ref = &finalizerRef{
		parent: f,
	}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
	return f
}

func init() {
	_ = NewTicker()
}
