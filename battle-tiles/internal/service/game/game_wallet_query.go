// internal/service/game/wallet_query_service.go
package game

import (
	gameBiz "battle-tiles/internal/biz/game"
	"battle-tiles/internal/dal/req"
	resp "battle-tiles/internal/dal/resp"
	plazaHTTP "battle-tiles/internal/utils/plaza"
	"battle-tiles/pkg/plugin/middleware"
	"battle-tiles/pkg/utils/ecode"
	"battle-tiles/pkg/utils/response"
	"time"

	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type WalletQueryService struct{ uc *gameBiz.FundsUseCase }

func NewWalletQueryService(uc *gameBiz.FundsUseCase) *WalletQueryService {
	return &WalletQueryService{uc: uc}
}

func (s *WalletQueryService) RegisterRouter(r *gin.RouterGroup) {
	g := r.Group("/members").Use(middleware.JWTAuth())
	g.POST("/wallet/get", middleware.RequirePerm("fund:wallet:view"), s.Get)     // 单人余额
	g.POST("/wallet/list", middleware.RequirePerm("fund:wallet:view"), s.List)   // 批量筛选
	g.POST("/ledger/list", middleware.RequirePerm("fund:ledger:view"), s.Ledger) // 流水
	// 用户战绩明细（基于外部HTTP数据源）
	g.POST("/battle/details", middleware.RequirePerm("battle:detail:view"), s.BattleDetails)
	// 本地战绩（已落库）
	g.POST("/battle/export/html", middleware.RequirePerm("battle:detail:export"), s.ExportBattleHTML)
}

// Get
// @Summary      查询单人钱包
// @Tags         资金/钱包
// @Accept       json
// @Produce      json
// @Param        in body req.GetWalletRequest true "house_gid, member_id"
// @Success      200 {object} response.Body{data=resp.WalletVO}
// @Router       /members/wallet/get [post]
func (s *WalletQueryService) Get(c *gin.Context) {
	var in req.GetWalletRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	m, err := s.uc.GetWallet(c.Request.Context(), in.HouseGID, in.MemberID)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, m)
}

// List
// @Summary      钱包列表（支持余额区间/是否个性额度过滤）
// @Tags         资金/钱包
// @Accept       json
// @Produce      json
// @Param        in body req.ListWalletsRequest true "house_gid, min/max_balance, has_custom_limit, page"
// @Success      200 {object} response.Body{data=resp.WalletListResponse}
// @Router       /members/wallet/list [post]
func (s *WalletQueryService) List(c *gin.Context) {
	var in req.ListWalletsRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	list, total, err := s.uc.ListWallets(c.Request.Context(), in.HouseGID, in.MinBalance, in.MaxBalance, in.HasCustomLimit, in.Page, in.PageSize)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.WalletListResponse{
		List:     list,
		Total:    total,
		Page:     normPage(in.Page),
		PageSize: normSize(in.PageSize),
	})
}

// Ledger
// @Summary      资金流水查询
// @Tags         资金/流水
// @Accept       json
// @Produce      json
// @Param        in body req.ListLedgerRequest true "house_gid, member_id 可选, type 可选, 时间范围"
// @Success      200 {object} response.Body{data=resp.LedgerListResponse}
// @Router       /members/ledger/list [post]
func (s *WalletQueryService) Ledger(c *gin.Context) {
	var in req.ListLedgerRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	var startPtr, endPtr *time.Time
	parse := func(s *string) *time.Time {
		if s == nil || *s == "" {
			return nil
		}
		// 尝试 RFC3339，再尝试 yyyy-mm-dd
		if t, err := time.Parse(time.RFC3339, *s); err == nil {
			return &t
		}
		if t, err := time.Parse("2006-01-02", *s); err == nil {
			tt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
			return &tt
		}
		return nil
	}
	startPtr = parse(in.StartAt)
	endPtr = parse(in.EndAt)
	// 若只给了 EndAt 的日期，默认 end = 次日 00:00:00
	if endPtr != nil && endPtr.Hour() == 0 && endPtr.Minute() == 0 && endPtr.Second() == 0 && in.EndAt != nil && len(*in.EndAt) == len("2006-01-02") {
		e := endPtr.AddDate(0, 0, 1)
		endPtr = &e
	}

	list, total, err := s.uc.ListLedger(c.Request.Context(), in.HouseGID, in.MemberID, in.Type, startPtr, endPtr, in.Page, in.PageSize)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	response.Success(c, resp.LedgerListResponse{
		List:     list,
		Total:    total,
		Page:     normPage(in.Page),
		PageSize: normSize(in.PageSize),
	})
}

func normPage(p int32) int32 {
	if p <= 0 {
		return 1
	}
	return p
}
func normSize(s int32) int32 {
	if s <= 0 || s > 200 {
		return 20
	}
	return s
}

// BattleDetails
// @Summary      用户战绩明细（外部HTTP）
// @Tags         战绩
// @Accept       json
// @Produce      json
// @Param        in body req.BattleDetailRequest true "house_gid, group_id, period(=today|yesterday|thisweek), 可选game_id"
// @Success      200 {object} response.Body{data=[]resp.BattleRecordVO}
// @Router       /members/battle/details [post]
func (s *WalletQueryService) BattleDetails(c *gin.Context) {
	var in req.BattleDetailRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}

	// period -> typeid (0=今日,1=昨日,2=本周)
	typeid := 0
	switch in.Period {
	case "today":
		typeid = 0
	case "yesterday":
		typeid = 1
	case "thisweek":
		typeid = 2
	default:
		typeid = 0
	}

	httpc := &httpClient{timeout: 10 * time.Second}
	list, err := plazaHTTP.GetGroupBattleInfoCtx(c.Request.Context(), httpc, "http://phone.foxuc.com/Ashx/GroService.ashx", in.GroupID, typeid)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}

	// 过滤到成员（可选）
	var out []resp.BattleRecordVO
	for _, b := range list {
		rec := resp.BattleRecordVO{RoomID: b.RoomID, KindID: b.KindID, BaseScore: b.BaseScore, Time: b.CreateTime}
		for _, p := range b.Players {
			if in.GameID != nil && *in.GameID != p.UserGameID {
				continue
			}
			rec.Players = append(rec.Players, resp.BattlePlayerVO{GameID: p.UserGameID, Score: p.Score})
		}
		if in.GameID == nil || len(rec.Players) > 0 {
			out = append(out, rec)
		}
	}
	response.Success(c, out)
}

// ExportBattleHTML 导出最近战绩为简单HTML（演示）
func (s *WalletQueryService) ExportBattleHTML(c *gin.Context) {
	// 简化：沿用在线查询后临时渲染为HTML（可改为查询本地库）
	var in req.BattleDetailRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, ecode.ParamsFailed, err)
		return
	}
	typeid := 0
	switch in.Period {
	case "today":
		typeid = 0
	case "yesterday":
		typeid = 1
	case "thisweek":
		typeid = 2
	}
	httpc := &httpClient{timeout: 10 * time.Second}
	list, err := plazaHTTP.GetGroupBattleInfoCtx(c.Request.Context(), httpc, "http://phone.foxuc.com/Ashx/GroService.ashx", in.GroupID, typeid)
	if err != nil {
		response.Fail(c, ecode.Failed, err)
		return
	}
	// 组装简单HTML
	b := &strings.Builder{}
	b.WriteString("<html><body>")
	for _, rec := range list {
		b.WriteString("<div style='margin:8px 0;padding:6px;border:1px solid #eee'>")
		b.WriteString(fmt.Sprintf("房:%d kind:%d base:%d 时间:%d<ul>", rec.RoomID, rec.KindID, rec.BaseScore, rec.CreateTime))
		for _, p := range rec.Players {
			b.WriteString(fmt.Sprintf("<li>%d: %d</li>", p.UserGameID, p.Score))
		}
		b.WriteString("</ul></div>")
	}
	b.WriteString("</body></html>")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(b.String()))
}

// 轻量HTTP客户端
type httpClient struct{ timeout time.Duration }

func (h *httpClient) Do(r *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: h.timeout}
	return client.Do(r)
}
