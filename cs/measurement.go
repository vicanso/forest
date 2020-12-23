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

package cs

const (
	// MeasurementPerformance 应用性能统计
	MeasurementPerformance = "performance"
	// MeasurementHTTPRequest http request统计
	MeasurementHTTPRequest = "httpRequest"
	// MeasurementRedisStats redis性能统计
	MeasurementRedisStats = "redisStats"
	// MeasurementRedisError redis出错统计
	MeasurementRedisError = "redisError"
	// MeasurementHTTPStats http性能统计
	MeasurementHTTPStats = "httpStats"
	// MeasurementHTTPInstanceStats http instance统计
	MeasurementHTTPInstanceStats = "httpInstanceStats"
	// MeasurementEntStats ent性能统计
	MeasurementEntStats = "entStats"
	// MeasurementEntOP ent的操作记录
	MeasurementEntOP = "entOP"
	// MeasurementHTTPError http响应出错统计
	MeasurementHTTPError = "httpError"
	// MeasurementUserTracker 用户行为记录
	MeasurementUserTracker = "userTracker"
	// MeasurementUserAction 用户行为记录
	// 用于前端记录客户相关的操作，如点击、确认、取消等
	MeasurementUserAction = "userAction"
	// MeasurementUserLogin 用户登录
	MeasurementUserLogin = "userLogin"
	// MeasurementUserAddTrack 添加用户跟踪
	MeasurementUserAddTrack = "userAddTrack"
	// MeasurementException 异常
	MeasurementException = "exception"
)
