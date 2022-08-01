package user_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

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

func Test_Update(t *testing.T) {
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
		Update(gomock.Any(), usr).
		Return(usr, nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserUpdated(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc)

	gotUser, err := svc.Update(context.TODO(), usr)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// TODO: improve asserts
	if gotUser == nil {
		t.Fatal("got nil user")
	}
}

func Test_Delete(t *testing.T) {
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
		Delete(gomock.Any(), usr).
		Return(nil)

	eventSvc := mock.NewEventService(ctrl)
	eventSvc.EXPECT().
		UserDeleted(gomock.Any(), usr).
		Return(nil)

	svc := user.NewService(mockStorage, eventSvc)

	err := svc.Delete(context.TODO(), usr)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
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
				// TODO: encode password
				// Password  "passwd"
				Email:   "email",
				Country: "DE",
			},
		},
		Total: 1,
	}

	opts := &user.ListOptions{Country: "DE", Search: "nick"}

	mockStorage := mock.NewStorage(ctrl)
	mockStorage.EXPECT().
		List(gomock.Any(), opts).
		Return(l, nil)

	svc := user.NewService(mockStorage, mock.NewEventService(ctrl))

	gotList, err := svc.List(context.TODO(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if gotList.Total != 1 {
		t.Fatalf("total = %d, want 1", gotList.Total)
	}
}
