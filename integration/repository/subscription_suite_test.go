package repository

import (
	"context"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/integration/helpers"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/models"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/repository"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SubscriptionRepositoryTestSuite struct {
	suite.Suite
	pgContainer *helpers.PostgresContainer
	repository  *repository.Repository
	ctx         context.Context
	db          *sqlx.DB
}

func (suite *SubscriptionRepositoryTestSuite) SetupSuite() {
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

	repo := repository.NewRepository(suite.db)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repo
}

func (suite *SubscriptionRepositoryTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %suite", err)
	}
}

func (suite *SubscriptionRepositoryTestSuite) TearDownTest() {
	err := suite.db.QueryRow("DELETE FROM subscription").Err()
	if err != nil {
		log.Fatalf("error clearing subscriptions table: %v", err)
	}
}

func (suite *SubscriptionRepositoryTestSuite) TestGetAllSubscriptions_Success() {
	t := suite.T()
	err := suite.repository.Subscription.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)

	subscriptions, err := suite.repository.Subscription.GetAllSubscriptions()
	assert.NoError(t, err)
	assert.NotNil(t, subscriptions)
	assert.Equal(t, subscriptions[0].Email, "example@gmail.com")
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_Success() {
	t := suite.T()
	err := suite.repository.Subscription.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_EmailAlreadyExists() {
	t := suite.T()
	err := suite.repository.Subscription.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
	err = suite.repository.Subscription.CreateSubscription("example@gmail.com")
	assert.ErrorIs(t, err, models.ErrEmailAlreadyExists)
}

func (suite *SubscriptionRepositoryTestSuite) TestDeleteSubscription_Success() {
	t := suite.T()
	err := suite.repository.Subscription.CreateSubscription("example@gmail.com")
	assert.NoError(t, err)
	err = suite.repository.Subscription.DeleteSubscription("example@gmail.com")
	assert.NoError(t, err)
	subscriptions, err := suite.repository.Subscription.GetAllSubscriptions()
	assert.NoError(t, err)
	assert.Empty(t, subscriptions)
}

func TestSubscriptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionRepositoryTestSuite))
}
