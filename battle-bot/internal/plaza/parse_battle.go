package plaza

import (
	"battle-bot/internal/plaza/game"
	"regexp"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var rePlusNum = regexp.MustCompile(`:(\s*)\+(\d+)`)

func tryInt(v any) (int, bool) {
	switch x := v.(type) {
	case float64:
		return int(x), true
	case string:
		s := strings.TrimSpace(strings.TrimPrefix(x, "+"))
		i, err := strconv.Atoi(s)
		return i, err == nil
	default:
		return 0, false
	}
}

func parseBattleInfo(body string) []*game.BattleInfo {
	tmp := make(map[string]any)
	if err := jsoniter.UnmarshalFromString(body, &tmp); err != nil {
		logger.Errorf("unmarshal battle info failed: %v", err)
		return nil
	}
	data, ok := tmp["Data"].(map[string]any)
	if !ok {
		return []*game.BattleInfo{}
	}
	list, ok := data["list"].([]any)
	if !ok {
		return []*game.BattleInfo{}
	}

	res := make([]*game.BattleInfo, 0, len(list))
	for _, item := range list {
		v, ok := item.(map[string]any)
		if !ok {
			continue
		}
		bi := &game.BattleInfo{}
		if i, ok := tryInt(v["MappedNum"]); ok {
			bi.RoomID = i
		}
		if i, ok := tryInt(v["KindID"]); ok {
			bi.KindID = i
		}
		if i, ok := tryInt(v["CreateTime"]); ok {
			bi.CreateTime = i
		}
		if i, ok := tryInt(v["BaseScore"]); ok {
			bi.BaseScore = i
		}
		raw, _ := v["UserScoreString"].(string)
		clean := rePlusNum.ReplaceAllString(raw, `:$1$2`)
		var dd []any
		if err := jsoniter.UnmarshalFromString(clean, &dd); err != nil {
			logger.Errorf("parse UserScoreString failed: %v, raw=%s", err, raw)
		} else {
			for _, d := range dd {
				m, _ := d.(map[string]any)
				var gid, score int
				if i, ok := tryInt(m["GameID"]); ok {
					gid = i
				}
				if i, ok := tryInt(m["Score"]); ok {
					score = i
				}
				bi.Players = append(bi.Players, &game.BattleSettle{
					UserGameID: gid,
					Score:      score,
				})
			}
		}
		res = append(res, bi)
	}
	return res
}
