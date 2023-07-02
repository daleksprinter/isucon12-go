package isuports

import (
	"context"
	"fmt"
	"sync"
)

type PlayerScoreCache struct {
	psc        map[string][]PlayerScoreRowWithPlayer
	updatedmap map[string]bool
	mu         sync.Mutex
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
		pss, err := GetPlayerScore(tenantDB, ctx, c.TenantID, c.ID)
		if err != nil {
			return err
		}
		d.update(c.ID, pss)
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

func (d *PlayerScoreCache) Get(tenantDB dbOrTx, ctx context.Context, tenantID int64, competitionID string) (pss []PlayerScoreRowWithPlayer, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if playerScoreCache.IsCached(competitionID) {
		pss, _ = playerScoreCache.psc[competitionID]
	} else {
		pss, err = GetPlayerScore(tenantDB, ctx, tenantID, competitionID)
		if err != nil {
			return nil, fmt.Errorf("error get player score%w", err)
		}
		playerScoreCache.update(competitionID, pss)
	}
	return pss, nil
}

func (d *PlayerScoreCache) update(competitionId string, data []PlayerScoreRowWithPlayer) {
	d.psc[competitionId] = data
	d.updatedmap[competitionId] = false
}

func (d *PlayerScoreCache) Updated(competitionid string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.updatedmap[competitionid] = true
}
