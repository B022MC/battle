package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	plazautils "battle-tiles/internal/utils/plaza"
)

func main() {
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}
	baseURL := "http://phone2.foxuc.com/Ashx/GroService.ashx"
	houseGID := 60870

	// 测试不同的时间范围
	typeids := []struct {
		id   int
		desc string
	}{
		{0, "最近3分钟"},
		{1, "最近30分钟"},
		{2, "最近1小时"},
	}

	for _, t := range typeids {
		fmt.Printf("\n=== 测试 %s (typeid=%d) ===\n", t.desc, t.id)
		battles, err := plazautils.GetGroupNewBattleInfoCtx(ctx, httpClient, baseURL, houseGID, t.id)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}
		
		fmt.Printf("获取到 %d 条战绩记录\n", len(battles))
		
		if len(battles) > 0 {
			fmt.Printf("\n前3条记录:\n")
			for i, b := range battles {
				if i >= 3 {
					break
				}
				fmt.Printf("  [%d] RoomID=%d, KindID=%d, BaseScore=%d, CreateTime=%d, Players=%d\n",
					i+1, b.RoomID, b.KindID, b.BaseScore, b.CreateTime, len(b.Players))
				for j, p := range b.Players {
					fmt.Printf("      Player[%d]: GameID=%d, Score=%d\n", j+1, p.UserGameID, p.Score)
				}
			}
		}
	}
}

