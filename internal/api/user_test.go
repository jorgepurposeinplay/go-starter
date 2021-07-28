package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oakeshq/go-starter/config"
	"github.com/oakeshq/go-starter/internal/storage"
	"github.com/oakeshq/go-starter/pkg/router"
	"github.com/oakeshq/go-starter/pkg/storage/dialer"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UserTestSuite struct {
	suite.Suite
	api *API
}

func TestCollection(t *testing.T) {
	cfg, err := config.LoadTestConfig("../../config/config.test.yaml")
	if err != nil {
		zap.L().Sugar().Fatalf("Unable to read config %v", err)
	}

	db, err := dialer.Connect(cfg)
	require.NoError(t, err)

	r := router.NewRouter()

	// RegisterHandlers(r, db, cfg)
	api := NewAPI(cfg, r, db)

	ts := &UserTestSuite{
		api: api,
	}

	suite.Run(t, ts)
}

func (ts *UserTestSuite) SetupTest() {
	ts.api.db.Exec("TRUNCATE TABLE users")
}

func (ts *UserTestSuite) TestCollection_UserList() {
	req, _ := http.NewRequest("GET", "/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")

	user := &storage.User{
		Email: "my@email.com",
	}

	createUser := ts.api.db.Create(user)
	require.NoError(ts.T(), createUser.Error)

	response := httptest.NewRecorder()

	ts.api.r.ServeHTTP(response, req)

	expected := fmt.Sprintf(`[
		{
			"email": "my@email.com"
		}
	]`)

	body, _ := ioutil.ReadAll(response.Body)

	ts.Assert().Equal(response.Code, http.StatusOK)
	ts.Assert().JSONEq(string(body), expected)
}

func (ts *UserTestSuite) TestCollection_UserListEmpty() {
	req, _ := http.NewRequest("GET", "/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	ts.api.r.ServeHTTP(response, req)

	expected := `[]`

	body, _ := ioutil.ReadAll(response.Body)

	ts.Assert().Equal(response.Code, http.StatusOK)
	ts.Assert().JSONEq(string(body), expected)
}
