package stat

import (
	"time"
	"url_shortener/pkg/db"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	var stat Stat
	currentDate := datatypes.Date(time.Now())

	result := repo.DB.Where("link_id=? and date=?", linkId, currentDate).First(&stat)
	// fmt.Println(stat)
	if result.Error != nil {
		repo.DB.Create(&Stat{
			LinkId: linkId,
			Clicks: 1,
			Date:   currentDate,
		})
	} else {
		stat.Clicks += 1
		repo.Save(&stat)
	}
}
func (repo *StatRepository) GetStats(from, to time.Time, by string) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case "month":
		selectQuery = "to_char(date,'YYYY-MM') as period, sum(clicks)"
	case "day":
		selectQuery = "to_char(date,'YYYY-MM-DD') as period, sum(clicks)"
	}

	repo.DB.Table("stats").
		Select(selectQuery).
		Where("date between ? and ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)
	return stats
}
