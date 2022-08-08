package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/mock"
)

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usr := &user.User{
		FirstName: "first",
		LastName:  "last",
		Nickname:  "nick",
		Password:  "passwd",
		Email:     "email",
		Country:   "DE",
	}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: usr.Email}).
		Return(&user.List{}, nil)

	mockStorage.EXPECT().
		Save(gomock.Any(), usr).
		Return(usr, nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserCreated(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc, 5)

	gotUser, err := svc.Save(context.TODO(), usr)
	require.NoError(t, err)
	require.Equal(t, usr, gotUser)
}

func Test_Create_AlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usr := &user.User{
		FirstName: "first",
		LastName:  "last",
		Nickname:  "nick",
		Password:  "passwd",
		Email:     "email",
		Country:   "DE",
	}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		List(gomock.Any(), &user.ListOptions{Search: usr.Email}).
		Return(&user.List{
			Total: uint64(1),
			Users: []*user.User{usr},
		}, nil)

	eventSvc := mock.NewEventService(ctrl)

	svc := user.NewService(mockStorage, eventSvc, 5)

	gotUser, err := svc.Save(context.TODO(), usr)
	require.ErrorIs(t, err, user.ErrAlreadyExists)
	require.Nil(t, gotUser)
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usr := &user.User{
		FirstName: "first",
		LastName:  "last",
		Nickname:  "nick",
		Password:  "passwd",
		Email:     "email",
		Country:   "DE",
	}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		Save(gomock.Any(), usr).
		Return(usr, nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserUpdated(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc, 5)

	gotUser, err := svc.Update(context.TODO(), usr)
	require.NoError(t, err)
	require.Equal(t, usr, gotUser)
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usr := &user.User{
		FirstName: "first",
		LastName:  "last",
		Nickname:  "nick",
		Password:  "passwd", Email: "email",
		Country: "DE",
	}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		Delete(gomock.Any(), usr).
		Return(nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserDeleted(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc, 5)

	err := svc.Delete(context.TODO(), usr)
	require.NoError(t, err)
}

func Test_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := &user.List{
		Users: []*user.User{
			{
				FirstName: "first",
				LastName:  "last",
				Nickname:  "nick",
				Email:     "email",
				Password:  "passwd",
				Country:   "DE",
			},
		},
		Total: 1,
	}

	opts := &user.ListOptions{Country: "DE", Search: "nick"}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		List(gomock.Any(), opts).
		Return(l, nil)

	svc := user.NewService(mockStorage, mock.NewEventService(ctrl), 5)

	gotList, err := svc.List(context.TODO(), opts)
	require.NoError(t, err)
	require.Equal(t, uint64(1), gotList.Total)
}
