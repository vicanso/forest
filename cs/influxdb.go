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
	// MeasurementRedisOP redis操作
	MeasurementRedisOP = "redisOP"
	// MeasurementRedisError redis出错统计
	MeasurementRedisError = "redisError"
	// MeasurementRedisConn redis连接
	MeasurementRedisConn = "redisConn"
	// MeasurementRouterConcurrency 路由并发访问
	MeasurementRouterConcurrency = "routerConcurrency"
	// MeasurementHTTPStats http性能统计
	MeasurementHTTPStats = "httpStats"
	// MeasurementHTTPInstanceStats http instance统计
	MeasurementHTTPInstanceStats = "httpInstanceStats"
	// MeasurementEntStats ent性能统计
	MeasurementEntStats = "entStats"
	// MeasurementEntUpdate ent的更新记录
	MeasurementEntUpdate = "entUpdate"
	// MeasurementEntQuery ent的查询记录
	MeasurementEntQuery = "entQuery"
	// MeasurementInfluxdbStats influxdb统计
	MeasurementInfluxdbStats = "influxdbStats"
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
	// MeasurementEvent 事件
	MeasurementEvent = "event"
)

const (
	// TagCategory 分类
	TagCategory = "category"
	// TagRoute 路由
	TagRoute = "route"
	// TagService 服务名称
	TagService = "service"
	// TagAction 用户的操作action
	TagAction = "action"
	// TagStep 用户操作的步骤
	TagStep = "step"
	// TagResult 操作结果
	// 由于influxdb中会使用result表示查询结果，因此如果使用result会导致展示时多了一个属性，
	// 因此调整名称
	TagResult = "rslt"
	// TagSchema 数据库的schema
	TagSchema = "schema"
	// TagOP 数据库的操作
	TagOP = "op"
	// TagMethod http method
	TagMethod = "method"
)

// string 类型
const (
	// FieldIP ip
	FieldIP = "ip"
	// FieldAddr addr
	FieldAddr = "addr"
	// FieldURI uri
	FieldURI = "uri"
	// FieldRouteName routeName
	FieldRouteName = "routeName"
	// FieldPath path
	FieldPath = "path"
	// FieldAccount 账号
	FieldAccount = "account"
	// FieldSID session id
	FieldSID = "sid"
	// FieldTID track id
	FieldTID = "tid"
	// FieldQuery url query
	FieldQuery = "query"
	// FieldParams url route params
	FieldParams = "params"
	// FieldForm request body
	FieldForm = "form"
	// FieldError error message
	FieldError = "error"
	// FieldUserAgent user agent
	FieldUserAgent = "userAgent"
	// FieldCountry 国家
	FieldCountry = "country"
	// FieldProvince 省份
	FieldProvince = "province"
	// FieldCity 城市
	FieldCity = "city"
	// FieldISP ISP
	FieldISP = "isp"
	// FieldErrCategory 出错分类
	FieldErrCategory = "errCategory"
)

// int 类型
const (
	// FieldMaxConcurrency 限制的最大并发数
	FieldMaxConcurrency = "maxConcurrency"
	// FieldProcessing 正在处理请求数
	FieldProcessing = "processing"
	// FieldTotalProcessing 正在处理的总请求数
	FieldTotalProcessing = "totalProcessing"
	// FilePipeProcessing pipe的正在处理请求数
	FilePipeProcessing = "pipeProcessing"
	// FieldLatency 耗时
	FieldLatency = "latency"
	// FieldStatus 状态码
	FieldStatus = "status"
	// FieldDNSUse dns耗时
	FieldDNSUse = "dnsUse"
	// FieldTCPUse tcp耗时
	FieldTCPUse = "tcpUse"
	// FieldTLSUse tls耗时
	FieldTLSUse = "tlsUse"
	// FieldProcessingUse 服务器处理耗时
	FieldProcessingUse = "processingUse"
	// FieldTransferUse 数据传输耗时
	FieldTransferUse = "transferUse"
	// FieldCount 总数
	FieldCount = "count"
	// FieldSize 大小
	FieldSize = "size"
	// FieldBodySize 内容大小
	FieldBodySize = "bodySize"
	// FieldHits 命中数量
	FieldHits = "hits"
	// FieldMisses miss数量
	FieldMisses = "misses"
	// FieldTimeouts 超时数量
	FieldTimeouts = "timeouts"
	// FieldTotalConns 总连接
	FieldTotalConns = "totalConns"
	// FieldIdleConns idle连接数
	FieldIdleConns = "idleConns"
	// FieldStaleConns stale连接数
	FieldStaleConns = "staleConns"
	// FieldMaxOpenConns 最大的连接数
	FieldMaxOpenConns = "maxOpenConns"
	// FieldOpenConns 当前连接数
	FieldOpenConns = "openConns"
	// FieldInUseConns 正在使用的连接
	FieldInUseConns = "inUseConns"
	// FieldWaitCount 等待的总数
	FieldWaitCount = "waitCount"
	// FieldWaitDuration 等待的时间
	FieldWaitDuration = "waitDuration"
	// FieldMaxIdleClosed idle close count
	FieldMaxIdleClosed = "maxIdleClosed"
	// FieldMaxIdleTimeClosed idle time close
	FieldMaxIdleTimeClosed = "maxIdleTimeClosed"
	// FieldMaxLifetimeClosed life time close
	FieldMaxLifetimeClosed = "maxLifetimeClosed"

	// FieldConnProcessing conn processing
	FieldConnProcessing = "connProcessing"
	// FieldConnAlive conn alive
	FieldConnAlive = "connAlive"
	// FieldConnCreatedCount conn created count
	FieldConnCreatedCount = "connCreatedCount"

	// num of fd
	FieldNumFds = "numFds"

	// FieldTotal 总数
	FieldTotal = "total"
	// FieldPoolSize pool size
	FieldPoolSize = "poolSize"
)

// bool 类型
const (
	// FieldReused 是否复用
	FieldReused = "reused"
	// FieldException 是否异常
	FieldException = "exception"
)

// map[string]any 类型
const (
	// FieldData 数据
	FieldData = "data"
)
