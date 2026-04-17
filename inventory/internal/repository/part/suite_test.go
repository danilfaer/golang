package part

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/danilfaer/golang/inventory/internal/repository/mocks"
)

var integrationColl *mongo.Collection

type RepositorySuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
}

func (s *RepositorySuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
}

func (s *RepositorySuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
