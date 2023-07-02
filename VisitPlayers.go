package isuports

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
)

type VisitersByCompetition struct {
	visiters map[string]map[string]bool
	mu       sync.Mutex
}

func (d *VisitersByCompetition) Initialize(t dbOrTx) error {
	p, err := d.getVisitorsByCompetitions(t, context.Background())
	if err != nil {
		return err
	}
	d.visiters = d.mappify(p)
	return nil
}

type VisitHistorySummaryRow struct {
	PlayerID      string `db:"player_id"`
	CompetitionID string `db:"competition_id"`
}

func (d *VisitersByCompetition) getVisitorsByCompetitions(t dbOrTx, ctx context.Context) ([]VisitHistorySummaryRow, error) {
	vhs := []VisitHistorySummaryRow{}
	if err := adminDB.SelectContext(
		ctx,
		&vhs,
		`SELECT player_id, competition_id 
				FROM visit_history vh join competition c on vh.competition_id = c.id 
				WHERE vh.created_at < c.finished_at 
				GROUP BY player_id, competition_id`,
	); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error Select visit_history: %w", err)
	}
	return vhs, nil
}

func (d *VisitersByCompetition) mappify(data []VisitHistorySummaryRow) map[string]map[string]bool {
	ret := map[string]map[string]bool{}
	for _, v := range data {
		if _, ok := ret[v.CompetitionID]; !ok {
			ret[v.CompetitionID] = map[string]bool{}
		}
		ret[v.CompetitionID][v.PlayerID] = true
	}
	return ret
}

func (d *VisitersByCompetition) GetVisitorsByCompetition(competitionID string) []string {
	visitorIDs := []string{}
	d.mu.Lock()
	defer d.mu.Unlock()
	p, ok := d.visiters[competitionID]
	if ok {
		for k, _ := range p {
			visitorIDs = append(visitorIDs, k)
		}
	}
	return visitorIDs
}

func (d *VisitersByCompetition) AddVisotor(comp CompetitionRow, v Viewer) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.visiters[comp.ID]; !ok {
		d.visiters[comp.ID] = map[string]bool{}
	}
	// competition.finished_atよりもあとの場合は、終了後に訪問したとみなして大会開催内アクセス済みとみなさない
	if comp.FinishedAt.Valid {
		return
	}
	d.visiters[comp.ID][v.playerID] = true
}
