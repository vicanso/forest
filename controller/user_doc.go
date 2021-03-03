// +build !codeanalysis
// 用户相关接口文档

package controller

// 用户信息
// swagger:response apiUserInfoResponse
type apiUserInfoResponse struct {
	// in: body
	Body *userInfoResp
}
