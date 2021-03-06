// +build !codeanalysis
// 用户相关接口文档

package controller

// 用户列表响应
//
// swagger:response apiUserListResponse
type apiUserListResponser struct {
	// in: body
	Body *userListResp
}

// 用户信息查询参数
//
// swagger:parameters userList
type apiUserListParams struct {
	userListParams
}

// 用户信息
// swagger:response apiUserInfoResponse
type apiUserInfoResponse struct {
	// in: body
	Body *userInfoResp
}
