package service

import (
	"errors"

	"github.com/siestacloud/service-monitoring/internal/core"
	"github.com/siestacloud/service-monitoring/internal/server/repository"
)

type RAMService struct {
	repo repository.RAM
}

// Конструктор - создает сервис для работы с мапкой метрик
func newRAMService(repo repository.RAM) *RAMService {
	return &RAMService{
		repo: repo,
	}
}

func (r *RAMService) CheckHash(key string, mtrx *core.Metric) error {

	hash := mtrx.GetHash()

	err := mtrx.SetHash(key)
	if err != nil {
		return err
	}
	nhash := mtrx.GetHash()
	if hash == nhash {
		return nil
	}
	return errors.New("hashes are not compared " + hash + " != " + nhash)
}

func (r *RAMService) GetAlljson() ([]byte, error) {
	return r.repo.GetAlljson()
}

func (r *RAMService) LookUP(key string) *core.Metric {
	return r.repo.LookUP(key)
}

func (r *RAMService) Update(key string, mtrx *core.Metric) error {
	return r.repo.Update(key, mtrx)
}

func (r *RAMService) Create(key string, mtrx *core.Metric) error {
	return r.repo.Create(key, mtrx)
}

// Добавляем новую метрику в мапу (обновляем существующую или создаем новую)
func (r *RAMService) Add(key string, mtrx *core.Metric) error {
	err := r.repo.Update(key, mtrx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RAMService) WriteLocalStorage(fn string) error {
	return r.repo.WriteLocalStorage(fn)
}

func (r *RAMService) ReadLocalStorage(fn string) (*core.MetricsPool, error) {
	return r.repo.ReadLocalStorage(fn)
}
