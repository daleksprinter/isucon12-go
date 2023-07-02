package isuports

import (
	"context"
	"fmt"
)

type PlayerScoreCache struct {
	psc        map[string][]PlayerScoreRowWithPlayer
	updatedmap map[string]bool
}

func (d *PlayerScoreCache) Initialize(tenantDB dbOrTx) error {
	d.psc = map[string][]PlayerScoreRowWithPlayer{}
	d.updatedmap = map[string]bool{}

	ctx := context.Background()
	competitions := []CompetitionRow{}
	if err := tenantDB.SelectContext(
		ctx,
		&competitions,
		`SELECT * FROM competition`,
	); err != nil {
		return fmt.Errorf("error Select player_score: %w", err)
	}
	for _, c := range competitions {
		pss := []PlayerScoreRowWithPlayer{}
		if err := tenantDB.SelectContext(
			ctx,
			&pss,
			`SELECT ps.score, ps.player_id, ps.row_num, p.display_name 
					FROM player_score_new ps join player p on ps.player_id = p.id 
					WHERE ps.tenant_id = ? and ps.competition_id = ? 
					`,
			c.TenantID,
			c.ID,
		); err != nil {
			return fmt.Errorf("error Select player_score:competitionID=%s, %w", c.ID, err)
		}
		d.Update(c.ID, pss)
	}
	return nil
}

func (d *PlayerScoreCache) IsCached(competitionid string) bool {
	ret, ok := d.updatedmap[competitionid]
	if !ok {
		return false
	}
	return !ret
}

func (d *PlayerScoreCache) Update(competitionId string, data []PlayerScoreRowWithPlayer) {
	d.psc[competitionId] = data
	d.updatedmap[competitionId] = false
}

func (d *PlayerScoreCache) Updated(competitionid string) {
	d.updatedmap[competitionid] = true
}
