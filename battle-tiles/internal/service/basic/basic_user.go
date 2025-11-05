package basic

import (
	basicBiz "battle-tiles/internal/biz/basic"
	rbacstore "battle-tiles/internal/dal/repo/rbac"
	"battle-tiles/internal/dal/req"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"

	_ "battle-tiles/internal/dal/model/basic"

	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
)

// BasicUserService is a basicUser service.
type BasicUserService struct {
	uc   *basicBiz.BasicUserUseCase
	rbac *rbacstore.Store
}

// NewBasicUserService new a basicUser service.
func NewBasicUserService(uc *basicBiz.BasicUserUseCase, rbac *rbacstore.Store) *BasicUserService {
	return &BasicUserService{uc: uc, rbac: rbac}
}

func (s *BasicUserService) RegisterRouter(rootRouter *gin.RouterGroup) {

	privateRouter := rootRouter.Group("/basic/user").Use(middleware.JWTAuth())
	privateRouter.POST("/addOne", s.AddOne)
	privateRouter.GET("/delOne", s.DelOne)
	privateRouter.POST("/delMany", s.DelMany)
	privateRouter.GET("/getOne", s.GetOne)
	privateRouter.GET("/getList", s.GetList)
	privateRouter.GET("/getOption", s.GetOption)
	privateRouter.POST("/updateOne", s.UpdateOne)

	// 新增：我的角色/权限查询
	privateRouter.GET("/me/roles", s.MeRoles)
	privateRouter.GET("/me/perms", s.MePerms)
	// 新增：我的平台账号
	privateRouter.GET("/me", s.Me)
	// 新增：修改我的登录密码
	privateRouter.POST("/changePassword", s.ChangePassword)
}

// ChangePasswordRequest 修改密码请求体
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// MeRoles 返回当前用户的角色ID（从 JWT claims）
// @Summary      我的角色ID
// @Tags         基础管理/用户
// @Security     BearerAuth
// @Produce      json
// @Success      200        {object}  response.Body
// @Router       /basic/user/me/roles [get]
func (s *BasicUserService) MeRoles(c *gin.Context) {
	cl, err := utils.GetClaims(c)
	if err != nil {
		response.Success(c, gin.H{"role_ids": []int32{}})
		return
	}
	response.Success(c, gin.H{"role_ids": cl.Roles})
}

// MePerms 返回当前用户的权限码（从 rbac store 聚合）
// @Summary      我的权限码
// @Tags         基础管理/用户
// @Security     BearerAuth
// @Produce      json
// @Success      200        {object}  response.Body
// @Router       /basic/user/me/perms [get]
func (s *BasicUserService) MePerms(c *gin.Context) {
	cl, err := utils.GetClaims(c)
	if err != nil {
		response.Success(c, gin.H{"perms": []string{}})
		return
	}
	set, err := s.rbac.GetUserPermCodes(c.Request.Context(), cl.UserID)
	if err != nil {
		response.Success(c, gin.H{"perms": []string{}})
		return
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	response.Success(c, gin.H{"perms": out})
}

// Me 返回当前用户的平台账号信息（从 JWT 解析 userID）
// @Summary      我的平台账号
// @Tags         基础管理/用户
// @Security     BearerAuth
// @Produce      json
// @Success      200        {object}  response.Body
// @Router       /basic/user/me [get]
func (s *BasicUserService) Me(c *gin.Context) {
	cl, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	one, err := s.uc.Get(c.Request.Context(), cl.UserID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, one)
}

// ChangePassword 修改当前登录用户的密码
// @Summary      修改密码
// @Tags         基础管理/用户
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        in body ChangePasswordRequest true "旧密码/新密码(明文)"
// @Success      200 {object} response.Body
// @Router       /basic/user/changePassword [post]
func (s *BasicUserService) ChangePassword(c *gin.Context) {
	var in ChangePasswordRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	claims, err := utils.GetClaims(c)
	if err != nil {
		response.Fail(c, ecode.TokenValidateFailed, err)
		return
	}
	me, err := s.uc.Get(c.Request.Context(), claims.BaseClaims.UserID)
	if err != nil || me == nil || me.BasicUser == nil {
		response.Fail(c, ecode.GetUserInfoFailed, "user not found")
		return
	}
	// 校验旧密码
	if !utils.BcryptCheck(in.OldPassword, me.Salt, me.Password) {
		response.Fail(c, ecode.UpdatePasswordFailed, "old password incorrect")
		return
	}
	if len(in.NewPassword) == 0 || len(in.NewPassword) > 30 {
		response.Fail(c, ecode.ParamsFailed, "invalid new_password length")
		return
	}
	salt, err := utils.GenerateSalt()
	if err != nil {
		response.Fail(c, ecode.UpdatePasswordFailed, err)
		return
	}
	hash := utils.BcryptHash(in.NewPassword, salt)
	if _, err := s.uc.Update(c.Request.Context(), claims.BaseClaims.UserID, map[string]interface{}{"password": hash, "salt": salt}); err != nil {
		response.Fail(c, ecode.UpdatePasswordFailed, err)
		return
	}
	response.SuccessWithOK(c)
}

// AddOne
// @Summary		    新增单条记录
// @Description	    新增单条记录 Add Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	body		req.AddBasicUserReq	true	"请求参数"
// @Success		    200		{object}	response.Body{data=basic.BasicUserDoc,msg=string}
// @Router			/basic/user/addOne [post]
func (s *BasicUserService) AddOne(ctx *gin.Context) {
	var inParam req.AddBasicUserReq
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	createdOne, err := s.uc.Create(ctx.Request.Context(), &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, createdOne)
}

// DelOne
// @Summary		    删除单条记录
// @Description	    删除单条记录 Del Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	query		utils.PkByInt32Param	true	"请求参数"
// @Success		    200		{object}	response.Body{msg=string}
// @Router			/basic/user/delOne [get]
func (s *BasicUserService) DelOne(ctx *gin.Context) {
	var inParam utils.PkByInt32Param
	if err := ctx.ShouldBindQuery(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Delete(ctx.Request.Context(), inParam.Id); err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(ctx)
}

// DelMany
// @Summary		    删除多条记录
// @Description	    删除多条记录 Del Many Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			id   query     []int32  true  "主键ID（多值：?id=1&id=2&id=3）"  collectionFormat(multi)
// @Success		    200		{object}	response.Body{msg=string}
// @Router			/basic/user/delMany [post]
func (s *BasicUserService) DelMany(ctx *gin.Context) {
	var inParam utils.PkByInt32sParam
	if err := ctx.ShouldBindQuery(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	if err := s.uc.Delete(ctx.Request.Context(), inParam.Id); err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.SuccessWithOK(ctx)
}

// GetOne
// @Summary		    查询单条记录
// @Description	    查询单条记录 By PK Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	query		utils.PkByInt32Param	true	"请求参数"
// @Success		    200		{object}	response.Body{data=basic.BasicUserDoc,msg=string}
// @Router			/basic/user/getOne [get]
func (s *BasicUserService) GetOne(ctx *gin.Context) {
	var inParam utils.PkByInt32Param
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	createdOne, err := s.uc.Get(ctx.Request.Context(), inParam.Id)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, createdOne)
}

// GetList
// @Summary		    查询N条记录
// @Description	    查询N条记录 List/Page Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	query		req.ListBasicUserReq	true	"请求参数"
// @Success		    200		{object}	response.Body{data=utils.PageResult{list=[]basic.BasicUserDoc},msg=string}
// @Router			/basic/user/getList [get]
func (s *BasicUserService) GetList(ctx *gin.Context) {
	var inParam req.ListBasicUserReq
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindQuery(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	list, total, err := s.uc.List(ctx.Request.Context(), &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, utils.PageResult{
		List:     list,
		NotPage:  inParam.NotPage,
		Total:    total,
		PageNo:   inParam.PageNo,
		PageSize: inParam.PageSize,
	})
}

// GetOption
// @Summary		    查询N条记录
// @Description	    查询N条记录 List/Page Model To Option
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	query		req.ListBasicUserReq	true	"请求参数"
// @Success		    200		{object}	response.Body{data=utils.PageResult{list=[]basic.BasicUserDoc},msg=string}
// @Router			/basic/user/getOption [get]
func (s *BasicUserService) GetOption(ctx *gin.Context) {
	var inParam req.ListBasicUserReq
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	list, total, err := s.uc.ListByOption(ctx.Request.Context(), &inParam)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, utils.PageResult{
		List:     list,
		NotPage:  inParam.NotPage,
		Total:    total,
		PageNo:   inParam.PageNo,
		PageSize: inParam.PageSize,
	})
}

// UpdateOne
// @Summary		    修改单条记录
// @Description	    修改单条记录 Update Model
// @Tags			基础管理/用户
// @Accept			json
// @Produce		    json
// @Param			inParam	body		req.AddBasicUserReq	true	"请求参数"
// @Success		    200		{object}	response.Body{data=basic.BasicUserDoc,msg=string}
// @Router			/basic/user/updateOne [post]
func (s *BasicUserService) UpdateOne(ctx *gin.Context) {
	var inParam req.UpdateBasicUserReq
	defaults.SetDefaults(&inParam)
	if err := ctx.ShouldBindJSON(&inParam); err != nil {
		response.Fail(ctx, ecode.ParamsFailed, err)
		return
	}
	// 解析token 中的用户信息 需要时使用
	_, err := utils.GetClaims(ctx)
	if err != nil {
		response.Fail(ctx, ecode.TokenValidateFailed, err)
		return
	}
	var fields map[string]interface{}
	mapstructure.Decode(inParam, &fields)
	updatedOne, err := s.uc.Update(ctx.Request.Context(), inParam.Id, fields)
	if err != nil {
		response.Fail(ctx, ecode.Failed, err)
		return
	}
	response.Success(ctx, updatedOne)
}
