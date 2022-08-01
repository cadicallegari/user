package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/cadicallegari/user"
	"github.com/cadicallegari/user/mock"
)

func Test_Create_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usr := &user.User{
		FirstName: "first",
		LastName:  "last",
		Nickname:  "nick",
		// TODO: encode password
		// Password  "passwd"
		Email:   "email",
		Country: "DE",
	}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		Create(gomock.Any(), usr).
		Return(usr, nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserCreated(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc)

	gotUser, err := svc.Create(context.TODO(), usr)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// TODO: improve asserts
	if gotUser == nil {
		t.Fatal("got nil user")
	}
}
