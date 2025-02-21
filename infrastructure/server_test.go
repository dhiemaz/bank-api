package infrastructure

import (
	"fmt"
	"github.com/dhiemaz/bank-api/config"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mockdb "github.com/dhiemaz/bank-api/infrastructure/db/mock"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testCase interface {
	buildStubsMethod(store *mockdb.MockStore)
	checkResponseMethod(t *testing.T, response *httptest.ResponseRecorder)
	setupAuthMethod(t *testing.T, request *http.Request, maker token.Maker)
}

type testCaseBase struct {
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, response *httptest.ResponseRecorder)
	setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
}

func (tcb testCaseBase) buildStubsMethod(store *mockdb.MockStore) { tcb.buildStubs(store) }
func (tcb testCaseBase) checkResponseMethod(t *testing.T, response *httptest.ResponseRecorder) {
	tcb.checkResponse(t, response)
}
func (tcb testCaseBase) setupAuthMethod(t *testing.T, request *http.Request, maker token.Maker) {
	tcb.setupAuth(t, request, maker)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func runServerTest(t *testing.T, tc testCase, req *http.Request) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	tc.buildStubsMethod(store)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	cfg := config.GetConfig()

	maker, err := token.NewPasetoMaker(cfg.SymmetricKey)
	if err != nil {
		panic(err)
	}

	port := "8080"
	tc.setupAuthMethod(t, req, maker)
	err = server.Start(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return
	}
	tc.checkResponseMethod(t, recorder)
}

func newTestServer(t *testing.T, store db.Store) *GinServer {
	testConfig := config.GetConfig()

	server, err := NewServer(testConfig, store, nil)
	require.NoError(t, err)
	return server
}
