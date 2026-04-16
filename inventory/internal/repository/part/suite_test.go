package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/danilfaer/golang/inventory/internal/repository/mocks"
)

type RepositorySuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
	repository     *repository
}

func (s *RepositorySuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
	// Здесь можно создать тестовую реализацию repository с моками
	// или тестировать только бизнес-логику без привязки к хранилищу
	s.repository = NewRepository()
}

func (s *RepositorySuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
