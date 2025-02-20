package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/dhiemaz/bank-api/infrastructure/db/mock"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/internal/handlers"
	"github.com/dhiemaz/bank-api/internal/handlers/account"
	"github.com/dhiemaz/bank-api/internal/handlers/user"
	"github.com/dhiemaz/bank-api/internal/middlewares"
	"github.com/dhiemaz/bank-api/util/token"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	user1, _ := user.createRandomUser(t)
	user2, _ := user.createRandomUser(t)

	account1 := account.createRandomAccount(user1.Username)
	account2 := account.createRandomAccount(user2.Username)

	amount := common.RandomInteger(1, 1000)
	account1.Currency, account2.Currency = common.IDR, common.IDR

	arg := createTransferReq{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	testCases := []struct {
		name          string
		transferArg   createTransferReq
		FromAccountID int64
		ToAccountID   int64
		handlers.testCaseBase
	}{
		{
			name:        "OK",
			transferArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.mockdb) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.db{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1)

					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account2, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					middlewares.addAuthHeader(t, req, maker, middlewares.authorizationTypeBearer, user1.Username)
				},
			},
		},
		{
			name: "BadRequest-Eq(IDS)",
			transferArg: createTransferReq{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
			},
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					middlewares.addAuthHeader(t, req, maker, middlewares.authorizationTypeBearer, user1.Username)
				},
			},
		},
		{
			name:        "BadRequest-Binding",
			transferArg: createTransferReq{},
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					middlewares.addAuthHeader(t, req, maker, middlewares.authorizationTypeBearer, user1.Username)
				},
			},
		},
		{
			name:        "InternalError",
			transferArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						TransferTx(gomock.Any(), gomock.Eq(db.TransferTxParam{
							FromAccountID: arg.FromAccountID,
							ToAccountID:   arg.ToAccountID,
							Amount:        arg.Amount,
						})).
						Times(1).Return(db.TransferTxResult{}, sql.ErrConnDone)

					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account1, nil)
					store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(account2, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
				setupAuth: func(t *testing.T, req *http.Request, maker token.Maker) {
					middlewares.addAuthHeader(t, req, maker, middlewares.authorizationTypeBearer, user1.Username)
				},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			data, err := json.Marshal(tc.transferArg)
			require.NoError(t, err)

			url := "/api/transfers"

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			handlers.runServerTest(t, tc, req)
		})
	}
}
