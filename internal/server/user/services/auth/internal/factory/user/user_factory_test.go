package user

import (
	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid/mock_uuid_generator"
	"github.com/anoriar/gophkeeper/internal/server/user/entity"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestUserFactory_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGeneratorMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)

	type args struct {
		login    string
		password string
		salt     string
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          entity.User
	}{
		{
			name: "success",
			args: args{
				login:    "test",
				password: "pass",
				salt:     "1234",
			},
			mockBehaviour: func() {
				uuidGeneratorMock.EXPECT().NewString().Return("ef9a2642-d0f7-4412-81c6-13d42ac88027")
			},
			want: entity.User{
				ID:       "ef9a2642-d0f7-4412-81c6-13d42ac88027",
				Login:    "test",
				Password: "pass",
				Salt:     "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			factory := NewUserFactory(uuidGeneratorMock)
			got := factory.Create(tt.args.login, tt.args.password, tt.args.salt)
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}
