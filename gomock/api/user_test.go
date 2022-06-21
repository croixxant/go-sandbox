package api_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/croixxant/go-sandbox/gomock/api"
	"github.com/croixxant/go-sandbox/gomock/db"
	mock_db "github.com/croixxant/go-sandbox/gomock/db/mock"
)

func randomUser() (db.User, string) {
	password := gofakeit.Password(true, true, true, false, false, 20)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return db.User{
		ID:             int64(gofakeit.Number(1, 1000)),
		Email:          gofakeit.Email(),
		HashedPassword: string(hashed),
		ConfirmedAt:    sql.NullTime{Valid: false},
		LikesCount:     int32(gofakeit.Number(1, 100)),
	}, password
}

func TestServer_getUser(t *testing.T) {
	user, _ := randomUser()

	tests := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
				assertBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: 0,
			buildStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := api.NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tt.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

			server.ServeHTTP(recorder, request)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestServer_createUser(t *testing.T) {
	user, password := randomUser()

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"Email":    user.Email,
				"Password": password,
			},
			buildStubs: func(store *mock_db.MockStore) {
				arg := db.CreateUserParams{
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tt.buildStubs(store)

			server := api.NewServer(store)
			recorder := httptest.NewRecorder()

			data, _ := json.Marshal(tt.body)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			assert.NoError(t, err)

			server.ServeHTTP(recorder, request)
			tt.checkResponse(recorder)
		})
	}
}

func assertBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	assert.NoError(t, err)

	var gotUser db.User
	_ = json.Unmarshal(data, &gotUser)
	assert.Equal(t, user, gotUser)
}

type eqCreateUserParamsMatcher struct {
	want     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	got, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(got.HashedPassword), []byte(e.password))
	if err != nil {
		return false
	}

	e.want.HashedPassword = got.HashedPassword
	return reflect.DeepEqual(e.want, got)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.want, e.password)
}

func EqCreateUserParams(want db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{want: want, password: password}
}
