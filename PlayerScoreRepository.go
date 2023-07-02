package isuports

import (
	"context"
	"fmt"
)

func GetPlayerScore(tenantDB dbOrTx, ctx context.Context, tenantID int64, competitionID string) ([]PlayerScoreRowWithPlayer, error) {
	pss := []PlayerScoreRowWithPlayer{}
	if err := tenantDB.SelectContext(
		ctx,
		&pss,
		`SELECT ps.score, ps.player_id, ps.row_num, p.display_name 
					FROM player_score_new ps join player p on ps.player_id = p.id 
					WHERE ps.tenant_id = ? AND ps.competition_id = ? 
					`,
		tenantID,
		competitionID,
	); err != nil {
		return nil, fmt.Errorf("error Select player_score: tenantID=%d, competitionID=%s, %w", tenantID, competitionID, err)
	}
	return pss, nil
}
