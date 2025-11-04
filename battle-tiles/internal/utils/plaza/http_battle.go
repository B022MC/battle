package plaza

import (
	"battle-tiles/internal/dal/vo/game"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

func GetGroupNewBattleInfoCtx(ctx context.Context, httpc HTTPDoer, base string, group, typid int) ([]*game.BattleInfo, error) {
	// base 例如: http://phone2.foxuc.com/Ashx/GroService.ashx
	st := strconv.FormatInt(time.Now().Unix(), 10)
	ps := url.Values{}
	ps.Set("groupid", strconv.Itoa(group))
	ps.Set("servertime", st)
	ps.Set("typeid", strconv.Itoa(typid))
	ps.Set("token", token(group, "WH3001", st))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"?action=GetGroupNewBattleInfo", strings.NewReader(ps.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoBattleClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errors.New("empty response")
	}
	return parseBattleInfo(string(body)), nil
}

func GetGroupBattleInfoCtx(ctx context.Context, httpc HTTPDoer, base string, group, typid int) ([]*game.BattleInfo, error) {
	// typeid=0 今日 / 1 昨日 / 2 本周
	st := strconv.FormatInt(time.Now().Unix(), 10)
	ps := url.Values{}
	ps.Set("groupid", strconv.Itoa(group))
	ps.Set("typeid", strconv.Itoa(typid))
	ps.Set("servertime", st)
	ps.Set("pageIndex", "1")
	ps.Set("PageSize", "50")
	ps.Set("StationID", "2000")
	ps.Set("token", token(group, "WH3001", st))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+"?action=getgroupbattleinfo", strings.NewReader(ps.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "GoBattleClient/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errors.New("empty response")
	}
	return parseBattleInfo(string(body)), nil
}
