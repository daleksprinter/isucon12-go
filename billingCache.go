package isuports

type Counts struct {
	Players  int64
	Visitors int64
}
type billingCache struct {
	memo map[string]Counts
}

func (d billingCache) Initialize() {
	d.memo = make(map[string]Counts)
}

func (d billingCache) update(competitionId string, count Counts) {
	d.memo[competitionId] = count
}

func (d billingCache) Get(competitionId string) Counts {
	if v, ok := d.memo[competitionId]; ok {
		return v
	} else {
		c := d.count(competitionId)
		d.update(competitionId, c)
		return c
	}
}

func (d billingCache) count(competitionId string) Counts {
	// ランキングにアクセスした参加者のIDを取得する
	billingMap := map[string]string{}
	visitors := visitorsByCompetition.GetVisitorsByCompetition(competitionId)
	for _, p := range visitors {
		billingMap[p] = "visitor"
	}

	// スコアを登録した参加者のIDを取得する
	for _, pid := range playersByCompetition.GetPlayersByCompetition(competitionId) {
		// スコアが登録されている参加者
		billingMap[pid] = "player"
	}

	// 大会が終了している場合のみ請求金額が確定するので計算する
	var playerCount, visitorCount int64
	for _, category := range billingMap {
		switch category {
		case "player":
			playerCount++
		case "visitor":
			visitorCount++
		}
	}
	return Counts{
		Players:  playerCount,
		Visitors: visitorCount,
	}
}
