package historyrepo

import (
	"github.com/DaniilKalts/music-platform-api/internal/adapter/database/postgres/sqlc"
	"github.com/DaniilKalts/music-platform-api/internal/domain/history"
)

func toDomainHistoryRecordFromList(row sqlc.ListListeningHistoryByUserIDRow) *history.HistoryRecord {
	return &history.HistoryRecord{
		ID:         row.ID,
		UserID:     uuidNil(),
		TrackID:    row.TrackID,
		TrackTitle: row.Title,
		ArtistName: row.ArtistName,
		ListenedAt: row.ListenedAt,
	}
}

func uuidNil() [16]byte {
	return [16]byte{}
}
