package basic

import (
	"battle-tiles/internal/biz/game"
	"battle-tiles/internal/conf"
	"battle-tiles/internal/consts"
	basicModel "battle-tiles/internal/dal/model/basic"
	basicRepo "battle-tiles/internal/dal/repo/basic"
	"battle-tiles/internal/dal/req"

	// basicVo "battle-tiles/internal/dal/vo/basic"
	resp "battle-tiles/internal/dal/resp"
	pdb "battle-tiles/pkg/plugin/dbx"
	"battle-tiles/pkg/utils"
	"battle-tiles/pkg/utils/request"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/wumansgy/goEncrypt/rsa"
)

type BasicLoginUseCase struct {
	repo          basicRepo.BasicLoginRepo
	authRepo      basicRepo.AuthRepo
	gameAccountUC *game.GameAccountUseCase // 游戏账号用例（可选，用于注册时绑定游戏账号）
	global        *conf.Global
	log           *log.Helper
}

func NewBasicLoginUseCase(repo basicRepo.BasicLoginRepo, global *conf.Global, authRepo basicRepo.AuthRepo, gameAccountUC *game.GameAccountUseCase, logger log.Logger) *BasicLoginUseCase {
	return &BasicLoginUseCase{
		repo:          repo,
		authRepo:      authRepo,
		gameAccountUC: gameAccountUC,
		global:        global,
		log:           log.NewHelper(log.With(logger, "module", "usecase/basic_login")),
	}
}

// 用户名 + 密码登录
func (uc *BasicLoginUseCase) LoginByUsernamePassword(ctx context.Context, c *gin.Context, req *req.UsernamePasswordLoginRequest) (*resp.LoginResponse, error) {
	user, err := uc.repo.FindByUsername(ctx, req.Username)
	if err != nil || user == nil {
		return nil, errors.New("用户不存在")
	}
	if user.Password == "" || user.Salt == "" {
		return nil, errors.New("该用户未设置密码")
	}

	passwordPlain, err := uc.BeforeValidatorPwd(ctx, req.Password)
	if err != nil {
		return nil, err
	}
	if !checkPassword(passwordPlain, user.Salt, user.Password) {
		return nil, errors.New("密码错误")
	}

	_ = uc.repo.UpdateLastLoginAt(ctx, user.Id)
	return uc.buildLoginResponse(c, user) // 改为调用方法，便于取权限
}

// 用户注册（用户名 + 密码 + 可选微信号 + 可选游戏账号）
func (uc *BasicLoginUseCase) Register(ctx context.Context, c *gin.Context, req *req.RegisterRequest) (*resp.LoginResponse, error) {
	if existUser, _ := uc.repo.FindByUsername(ctx, req.Username); existUser != nil {
		return nil, errors.New("用户名已存在")
	}

	passwordPlain, err := uc.BeforeValidatorPwd(ctx, req.Password)
	if err != nil {
		return nil, err
	}
	salt, err := utils.GenerateSalt()
	if err != nil {
		return nil, errors.New("生成密码 salt 失败")
	}
	passwordHash := utils.BcryptHash(passwordPlain, salt)

	username := req.Username
	if username == "" {
		// 理论上不会进来，binding 已经 required；保底逻辑保留
		for i := 0; i < 5; i++ {
			username = utils.GenerateUsername()
			exist, _ := uc.repo.FindByUsername(ctx, username)
			if exist == nil {
				break
			}
		}
	}
	nickName := req.NickName
	if nickName == "" {
		nickName = utils.GenerateNickName()
	}

	user := &basicModel.BasicUser{
		Username: username,
		Password: passwordHash,
		Salt:     salt,
		NickName: nickName,
		Avatar:   req.Avatar,
	}
	if req.WechatID != "" {
		user.WechatID = req.WechatID
	}

	id, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "注册失败")
	}
	user.Id = id

	// 为新注册用户绑定普通用户角色（忽略错误，不阻断注册）
	if uc.authRepo != nil {
		if err := uc.authRepo.EnsureUserHasOnlyRoleByCode(c.Request.Context(), user.Id, "ordinary-user"); err != nil {
			uc.log.Errorf("EnsureUserHasRoleByCode err: %v", err)
		}
	}

	// 如果提供了游戏账号信息，自动绑定游戏账号（忽略错误，不阻断注册）
	if req.GameAccount != "" && req.GamePassword != "" && req.GameAccountMode != "" && uc.gameAccountUC != nil {
		var mode consts.GameLoginMode
		switch req.GameAccountMode {
		case "account":
			mode = consts.GameLoginModeAccount
		case "mobile":
			mode = consts.GameLoginModeMobile
		default:
			uc.log.Warnf("invalid game account mode: %s", req.GameAccountMode)
			goto skipGameAccountBind
		}

		// 绑定游戏账号
		if _, err := uc.gameAccountUC.BindSingle(ctx, user.Id, mode, req.GameAccount, req.GamePassword, nickName); err != nil {
			uc.log.Errorf("bind game account failed: %v", err)
			// 不阻断注册流程，只记录错误
		} else {
			uc.log.Infof("user %d successfully bound game account", user.Id)
		}
	}

skipGameAccountBind:
	return uc.buildLoginResponse(c, user) // 新用户注册后同样返回携带权限的 token
}

// ===== 内部工具 =====

// 把平台、角色、权限写入 JWT，并返回登录响应
func (uc *BasicLoginUseCase) buildLoginResponse(c *gin.Context, user *basicModel.BasicUser) (*resp.LoginResponse, error) {
	j := utils.NewJWT()

	platform := ""
	if v := c.Request.Context().Value(pdb.CtxDBKey); v != nil {
		if s, ok := v.(string); ok {
			platform = s
		}
	}

	// 查询角色 & 权限（出错不阻断登录，仅记录日志）
	var roles []int32
	var perms []string
	if uc.authRepo != nil {
		var err error
		if roles, err = uc.authRepo.ListRoleIDsByUser(c.Request.Context(), user.Id); err != nil {
			uc.log.Errorf("ListRoleIDsByUser err: %v", err)
		}
		if perms, err = uc.authRepo.ListPermsByUser(c.Request.Context(), user.Id); err != nil {
			uc.log.Errorf("ListPermsByUser err: %v", err)
		}
	}

	claims := request.BaseClaims{
		UserID:   user.Id,
		Platform: platform,
		Username: user.Username,
		NickName: user.NickName,
		Roles:    roles,
		Perms:    perms,
	}
	tokenClaims := j.CreateClaims(claims)
	accessToken, refreshToken, err := j.CreateToken2(tokenClaims)
	if err != nil {
		return nil, err
	}

	return &resp.LoginResponse{
		User: &resp.BaseUserInfo{
			ID:           user.Id,
			Username:     user.Username,
			Avatar:       user.Avatar,
			NickName:     user.NickName,
			Introduction: user.Introduction,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(time.Until(tokenClaims.ExpiresAt.Time).Seconds()),
		Platform:     platform,
		Role:         user.Role, // 添加用户角色
		Roles:        roles,
		Perms:        perms,
	}, nil
}

func (uc *BasicLoginUseCase) BeforeValidatorPwd(ctx context.Context, password string) (string, error) {
	// 前端 RSA 加密 -> 后端私钥解密
	pwd, err := rsa.RsaDecryptByBase64(password, uc.global.Rsa.Private)
	if err != nil {
		uc.log.Infof("decrypt password err: %v", err)
		return "", err
	}
	if len(pwd) > 30 {
		return "", errors.New("密码不能超过30个字符")
	}
	return string(pwd), nil
}

func checkPassword(input, salt, encrypted string) bool {
	return utils.BcryptCheck(input, salt, encrypted)
}
