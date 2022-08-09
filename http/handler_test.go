package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/cadicallegari/user"
	userHttp "github.com/cadicallegari/user/http"
	"github.com/cadicallegari/user/mock"
	"github.com/cadicallegari/user/pkg/xhttp"
	"github.com/cadicallegari/user/pkg/xlogger"
)

type userTestSuite struct {
	ctx         context.Context
	log         *logrus.Entry
	router      *xhttp.Router
	svc         user.Service
	storageMock *mock.Storage
	eventMock   *mock.EventService
}

func serviceWithMocks(t *testing.T, ctrl *gomock.Controller) userTestSuite {
	var s userTestSuite

	s.log = xlogger.New(nil).WithFields(nil)
	s.ctx = xlogger.SetLogger(context.TODO(), s.log)
	s.storageMock = mock.NewStorage(ctrl)
	s.eventMock = mock.NewEventService(ctrl)

	s.svc = user.NewService(s.storageMock, s.eventMock, 4)

	s.router = xhttp.NewRouter(s.log)

	// to setup routes
	_ = userHttp.NewUserHandler(s.router, s.svc)

	return s
}

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	u := &user.User{
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: u.Email}).
		Return(&user.List{}, nil)

	suite.storageMock.EXPECT().
		Save(gomock.Any(), u).
		Return(u, nil)

	suite.eventMock.EXPECT().
		UserCreated(gomock.Any(), u).
		Return(nil)

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
}

func Test_Create_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	u := &user.User{
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: u.Email}).
		Return(&user.List{}, nil)

	suite.storageMock.EXPECT().
		Save(gomock.Any(), u).
		Return(nil, errors.New("any error"))

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Create_SendEventError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	u := &user.User{
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: u.Email}).
		Return(&user.List{}, nil)

	suite.storageMock.EXPECT().
		Save(gomock.Any(), u).
		Return(u, nil)

	suite.eventMock.EXPECT().
		UserCreated(gomock.Any(), u).
		Return(errors.New("some error"))

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	u := &user.User{
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: u.Email}).
		Return(&user.List{
			Total: uint64(1),
			Users: []*user.User{u},
		}, nil)

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func Test_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	id := "notfound"

	suite.storageMock.EXPECT().
		Get(gomock.Any(), id).
		Return(nil, user.ErrNotFound)

	req, err := http.NewRequest(http.MethodGet, "/v1/users/"+id, nil)
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	id := "some-valid-id"
	u := &user.User{
		ID:        id,
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		Get(gomock.Any(), id).
		Return(u, nil)

	suite.storageMock.EXPECT().
		Update(gomock.Any(), u).
		Return(u, nil)

	suite.eventMock.EXPECT().
		UserUpdated(gomock.Any(), u).
		Return(nil)

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/v1/users/"+id, bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Update_UserNotfound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	id := "notfound"
	u := &user.User{
		ID:        id,
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		Get(gomock.Any(), id).
		Return(nil, user.ErrNotFound)

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, "/v1/users/"+id, bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	id := "some-valid-id"
	u := &user.User{
		ID:        id,
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		Get(gomock.Any(), id).
		Return(u, nil)

	suite.storageMock.EXPECT().
		Delete(gomock.Any(), u).
		Return(nil)

	suite.eventMock.EXPECT().
		UserDeleted(gomock.Any(), u).
		Return(nil)

	req, err := http.NewRequest(http.MethodDelete, "/v1/users/"+id, nil)
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_Delete_UserNotfound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	suite := serviceWithMocks(t, ctrl)

	id := "notfound"
	u := &user.User{
		ID:        id,
		FirstName: "first name",
		Email:     "email",
	}

	suite.storageMock.EXPECT().
		Get(gomock.Any(), id).
		Return(nil, user.ErrNotFound)

	buf, err := json.Marshal(u)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodDelete, "/v1/users/"+id, bytes.NewBuffer(buf))
	require.NoError(t, err)

	req = req.WithContext(suite.ctx)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
