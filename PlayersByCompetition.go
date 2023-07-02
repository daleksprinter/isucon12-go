package isuports

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type PlayersByCompetition struct {
	players map[string]map[string]bool
	mu      sync.Mutex
}

func (d *PlayersByCompetition) Initialize(t dbOrTx) error {
	p, err := d.getPlayersByCompetitions(t, context.Background())
	if err != nil {
		return err
	}
	d.players = d.playersByCompetitionsToMap(p)
	return nil
}

func (d *PlayersByCompetition) getPlayersByCompetitions(tenantDB dbOrTx, ctx context.Context) ([]PlayerWithCompetition, error) {
	scoredPlayerIDs := []PlayerWithCompetition{}
	if err := tenantDB.SelectContext(
		ctx,
		&scoredPlayerIDs,
		"SELECT competition_id, player_id FROM player_score_new group by competition_id, player_id",
	); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error Select count player_score:  %w", err)
	}
	return scoredPlayerIDs, nil
}

func (d *PlayersByCompetition) playersByCompetitionsToMap(data []PlayerWithCompetition) map[string]map[string]bool {
	ret := map[string]map[string]bool{}
	for _, v := range data {
		if _, ok := ret[v.CompetitionID]; !ok {
			ret[v.CompetitionID] = map[string]bool{}
		}
		ret[v.CompetitionID][v.PlayerID] = true
	}
	return ret
}

func (d *PlayersByCompetition) GetPlayersByCompetition(competitionID string) []string {
	scoredPlayerIDs := []string{}
	d.mu.Lock()
	p, ok := d.players[competitionID]
	if ok {
		for k, _ := range p {
			scoredPlayerIDs = append(scoredPlayerIDs, k)
		}
	}
	d.mu.Unlock()
	return scoredPlayerIDs
}

func (d *PlayersByCompetition) AddPlayers(competitionID string, playerScoreRows []PlayerScoreRow) {
	d.mu.Lock()
	for _, v := range playerScoreRows {
		if _, ok := d.players[competitionID]; !ok {
			d.players[competitionID] = map[string]bool{}
		}
		d.players[v.CompetitionID][v.PlayerID] = true
	}
	d.mu.Unlock()
}
