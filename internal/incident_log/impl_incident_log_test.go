package incident_log

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/acidnyi/incident-log-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IncidentLogSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[IncidentLogDefinition]
}

func TestIncidentLogSuite(t *testing.T) {
	suite.Run(t, new(IncidentLogSuite))
}

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (this *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := this.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (this *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := this.Called(ctx, id, document)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := this.Called(ctx, id)
	return args.Error(0)
}

func (this *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := this.Called(ctx)
	return args.Error(0)
}

func (suite *IncidentLogSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[IncidentLogDefinition]{}

	// Compile time assert that the mock is of type db_service.DbService[IncidentLogDefinition]
	var _ db_service.DbService[IncidentLogDefinition] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&IncidentLogDefinition{
				Id:       "hospital-security-log",
				Name:     "Nemocničný bezpečnostný denník",
				Location: "Univerzitná nemocnica",
				Incidents: []Incident{
					{
						Id:           "INC-001",
						IncidentType: "Bezpečnostná udalosť",
						Location:     "Urgentný príjem",
						OccurredAt:   time.Now(),
						Description:  "Pôvodný popis",
					},
				},
			},
			nil,
		)

	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
}

func (suite *IncidentLogSuite) Test_UpdateIncident_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	json := `{
		"id": "INC-001",
		"incidentType": "Bezpečnostná udalosť",
		"location": "Urgentný príjem",
		"description": "Neoprávnený vstup do vyhradenej zóny.",
		"severity": "Vysoká",
		"status": "Nový"
	}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	ctx.Set("db_service", suite.dbServiceMock)

	ctx.Params = []gin.Param{
		{Key: "incidentLogId", Value: "hospital-security-log"},
		{Key: "incidentId", Value: "INC-001"},
	}

	ctx.Request = httptest.NewRequest(
		"PUT",
		"/incident-log/hospital-security-log/incidents/INC-001",
		strings.NewReader(json),
	)

	sut := implIncidentLogAPI{}

	// ACT
	sut.UpdateIncident(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(
		suite.T(),
		"UpdateDocument",
		mock.Anything,
		"hospital-security-log",
		mock.Anything,
	)
}
