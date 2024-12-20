package jobs

import (
	"github.com/kaellybot/kaelly-discord/models/constants"
	"github.com/kaellybot/kaelly-discord/models/entities"
	"github.com/kaellybot/kaelly-discord/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetJobs() ([]entities.Job, error) {
	var jobs []entities.Job
	response := repo.db.GetDB().
		Model(&entities.Job{}).
		Where("game = ?", constants.GetGame().AMQPGame).
		Preload("Labels").Find(&jobs)
	return jobs, response.Error
}
