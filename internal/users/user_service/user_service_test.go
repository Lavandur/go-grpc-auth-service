package user_service

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	roles_mock "auth-service/internal/roles/mock"
	users_mock "auth-service/internal/users/mock"
	"auth-service/pkg/logger"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_userService_Create(t *testing.T) {
	t.Parallel()

	password := "TEST_PASSWORD"
	expected := getUser(password)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepository := users_mock.NewMockUserRepository(ctrl)
	userRepository.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, user *models.User) (*models.User, error) {
			user.UserID = expected.UserID
			user.VisibleID = expected.VisibleID
			user.CreatedAt = expected.CreatedAt
			user.UpdatedAt = expected.UpdatedAt
			user.HashedPassword = expected.HashedPassword
			return user, nil
		}).
		AnyTimes()

	roleService := roles_mock.NewMockRoleService(ctrl)
	roleService.EXPECT().
		GetList(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(expected.Roles, nil).
		AnyTimes()

	log := logger.SetupLogger(nil)

	type args struct {
		ctx  context.Context
		data *models.UserInput
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Create user",
			args: args{
				ctx: context.Background(),
				data: &models.UserInput{
					Login:     expected.Login,
					Password:  password,
					Firstname: expected.Person.Firstname,
					Lastname:  expected.Person.Lastname,
					Birthdate: expected.Person.Birthdate,
					Email:     expected.Person.Email,
					Gender:    expected.Person.Gender,
					RoleIDs:   nil,
				},
			},
			want:    expected,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &userService{
				userRepository: userRepository,
				roleService:    roleService,
				logger:         log,
			}
			got, err := u.Create(tt.args.ctx, tt.args.data)
			tt.wantErr(t, err, "UserService.Create")
			assert.Equal(t, tt.want, got, "UserService.Create")
		})
	}
}

func getUser(password string) *models.User {
	dateTime := time.Now().UTC().Truncate(time.Microsecond)
	return &models.User{
		UserID:         uuid.New().String(),
		Login:          uuid.New().String(),
		VisibleID:      uuid.New().String(),
		HashedPassword: common.HashPassword(password),
		Person: models.Person{
			Firstname: "test",
			Lastname:  "test",
			Birthdate: dateTime,
			Email:     "test.test@gmail.com",
		},
		Roles: []*models.Role{{
			RoleID:      "SOME ROLE ID",
			Title:       "SOME ROLE NAME",
			Description: nil,
			CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
		}},
		CreatedAt:             dateTime,
		UpdatedAt:             dateTime,
		DeletedAt:             nil,
		LastPasswordRestoreAt: nil,
		SearchIndex:           nil,
	}
}
