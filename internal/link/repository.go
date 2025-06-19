package link

import (
	"url_shortener/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}
func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}
func (repo *LinkRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var l Link
	result := repo.Database.DB.Where("hash = ?", hash).Find(&l)
	if result.Error != nil {
		return nil, result.Error
	}
	return &l, nil
}

func (repo *LinkRepository) GetByID(id uint) (*Link, error) {
	var l Link
	result := repo.Database.DB.First(&l, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &l, nil
}

func (repo *LinkRepository) GetAll(limit, offset int) []Link {
	var links []Link

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(limit).
		Offset(offset).
		Scan(&links)

	return links
}
func (repo *LinkRepository) Count() int64 {
	var count int64

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Count(&count)

	return count
}
