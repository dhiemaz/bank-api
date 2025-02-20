package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/dhiemaz/bank-api/infrastructure/db/mock"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/internal/handlers"
	"github.com/dhiemaz/bank-api/util/token"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	user, password := createRandomUser(t)
	uniqueViolationError := &pq.Error{Code: "23505"}

	arg := createUserReq{
		Username:        user.Username,
		FullName:        user.FullName,
		Email:           user.Email,
		Password:        password,
		PasswordConfirm: password,
	}

	testCases := []struct {
		name    string
		userArg createUserReq
		handlers.testCaseBase
	}{
		{
			name:    "OK",
			userArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.mockdb) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(user, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusCreated, recorder.Code)
					requireBodyMatchUser(t, recorder.Body, user)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {},
			},
		},
		{
			name: "BadRequest",
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {},
			},
		},
		{
			name:    "SQLUserNameViolation",
			userArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, uniqueViolationError)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusForbidden, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {},
			},
		},
		{
			name:    "SQLEmailViolation",
			userArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, uniqueViolationError)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusForbidden, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {},
			},
		},
		{
			name:    "InternalErrorDB",
			userArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, sql.ErrConnDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]

		url := "/api/users/register"

		data, err := json.Marshal(tc.userArg)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
		require.NoError(t, err)

		handlers.runServerTest(t, tc, req)
	}

}

func requireBodyMatchUser(t *testing.T, b io.Reader, user db.User) {
	data, err := io.ReadAll(b)
	require.NoError(t, err)

	var userReceived db.User
	err = json.Unmarshal(data, &userReceived)
	require.NoError(t, err)
}

func createRandomUser(t *testing.T) (db.User, string) {
	password := common.RandomString(6)

	hashPassword, err := common.GenerateHashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Username:       common.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       common.RandomOwner(),
		Email:          common.RandomEmail(),
	}

	return user, password

}

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name string
		handlers.testCaseBase
	}{
		{
			name: "OK",
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {

				},
			},
		},
	}

	_ = testCases
}

func TestUpdateUser(t *testing.T) {
	var arg db.UpdateUserParams

	testCases := []struct {
		name    string
		userArg db.UpdateUserParams
		handlers.testCaseBase
	}{
		{
			name:    "OK",
			userArg: arg,
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						UpdateUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {

				},
			},
		},
	}

	_ = testCases
}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name string
		handlers.testCaseBase
	}{
		{
			name: "OK",
			testCaseBase: handlers.testCaseBase{
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetUser(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.User{}, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
				},
				setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				}},
		}}

	_ = testCases
}
