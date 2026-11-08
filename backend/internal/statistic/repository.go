package statistic

import (
	"context"
	"time"

	"github.com/exexeeeex/fiks-api/internal/domain"
	"gorm.io/gorm"
)

type StartTimeRepository interface {
	ListByGroupIds(
		ctx context.Context,
		groupIds []string,
		from, to time.Time,
	) ([]domain.StartTime, error)
}

type startTimeRepository struct {
	db *gorm.DB
}

func NewStartTimeRepository(db *gorm.DB) StartTimeRepository {
	return &startTimeRepository{db: db}
}

func (r *startTimeRepository) ListByGroupIds(
	ctx context.Context,
	groupIds []string,
	from, to time.Time,
) ([]domain.StartTime, error) {
	var startTimes []domain.StartTime
	if len(groupIds) == 0 {
		return startTimes, nil
	}

	err := r.db.WithContext(ctx).
		Where("group_id IN ? AND start_time >= ? AND start_time < ?", groupIds, from, to).
		Find(&startTimes).
		Error
	if err != nil {
		return nil, err
	}

	return startTimes, nil
}
