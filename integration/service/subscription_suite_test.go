package service

import (
	"context"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/integration/helpers"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/repository"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/service"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SubscriptionServiceTestSuite struct {
	suite.Suite
	pgContainer         *helpers.PostgresContainer
	subscriptionService *service.SubscriptionService
	ctx                 context.Context
	db                  *sqlx.DB
}

func (suite *SubscriptionServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := helpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer

	dbURL := pgContainer.ConnectionString

	suite.db, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewSubscriptionRepository(suite.db)
	suite.subscriptionService = service.NewSubscriptionService(repo)
}

func (suite *SubscriptionServiceTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %suite", err)
	}
}

func (suite *SubscriptionServiceTestSuite) TearDownTest() {
	err := suite.db.QueryRow("DELETE FROM subscription").Err()
	if err != nil {
		log.Fatalf("error clearing subscriptions table: %v", err)
	}
}

func (suite *SubscriptionServiceTestSuite) TestGetAllSubscriptions_Success() {
	t := suite.T()
	err := suite.subscriptionService.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)

	subscriptions, err := suite.subscriptionService.GetAllSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subscriptions)
	assert.Equal(t, subscriptions[0].Email, "example@gmail.com")
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_Success() {
	t := suite.T()
	err := suite.subscriptionService.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
}

func (suite *SubscriptionServiceTestSuite) TestCreateSubscription_EmailAlreadyExists() {
	t := suite.T()
	err := suite.subscriptionService.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
	err = suite.subscriptionService.CreateSubscription("example@gmail.com")
	assert.ErrorIs(t, err, models.ErrEmailAlreadyExists)
}

func (suite *SubscriptionServiceTestSuite) TestDeleteSubscription_Success() {
	t := suite.T()
	err := suite.subscriptionService.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
	err = suite.subscriptionService.DeleteSubscription("example@gmail.com")
	assert.NoError(t, err)
	subscriptions, err := suite.subscriptionService.GetAllSubscriptions()
	assert.NoError(t, err)
	assert.Empty(t, subscriptions)
}

func TestSubscriptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionServiceTestSuite))
}
