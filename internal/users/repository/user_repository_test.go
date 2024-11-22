package repository

import (
	"auth-service/internal/models"
	"auth-service/internal/roles"
	"auth-service/internal/roles/mock"
	"auth-service/testingdb"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func prepareUser() models.User {
	uid := "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"
	fTime := time.Now().UTC().Format(time.RFC850)
	someTime, _ := time.Parse(time.RFC850, fTime)
	return models.User{
		UserID:         uid,
		Login:          "70000000000",
		VisibleID:      "VI-SIBLE-ID",
		HashedPassword: "hashed_password",
		Person: models.Person{
			Firstname: "Alexey",
			Lastname:  "Somesamovich",
			Birthdate: someTime,
			Email:     "alexey.somesamovich@gmail.com",
		},
		Roles: []*models.Role{{
			RoleID:      "",
			Name:        "",
			Description: nil,
			CreatedAt:   time.Time{},
		}},
		CreatedAt:             someTime,
		UpdatedAt:             someTime,
		DeletedAt:             nil,
		LastPasswordRestoreAt: &someTime,
		SearchIndex:           nil,
	}
}

func getRoleRepMock(ctrl *gomock.Controller) roles.RoleRepository {
	roleRep := mock.NewMockRoleRepository(ctrl)
	roleRep.EXPECT().
		GetByID(gomock.Any(), gomock.Any()).
		Return(&models.Role{
			RoleID:      "",
			Name:        "",
			Description: nil,
			CreatedAt:   time.Time{},
		}, nil).AnyTimes()
	return roleRep
}

func Test_usersRepository_Create(t *testing.T) {
	t.Parallel()

	user := prepareUser()

	t.Run("Create user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)

		ctrl := gomock.NewController(t)
		roleRep := getRoleRepMock(ctrl)
		defer ctrl.Finish()

		repos := &usersRepository{pg.DB(), roleRep, nil}

		create, err := repos.Create(ctx, &user)
		if err != nil {
			assert.NoError(t, err)
		}
		reflect.DeepEqual(&user, create)
	})

	t.Run("Create user with same id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)

		ctrl := gomock.NewController(t)
		roleRep := getRoleRepMock(ctrl)
		defer ctrl.Finish()

		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Create(ctx, &user)
		if err != nil {
			assert.NoError(t, err)
		}
		_, err = repos.Create(ctx, &user)
		assert.Error(t, err)
	})
}

func Test_usersRepository_Delete(t *testing.T) {
	t.Parallel()

	user := prepareUser()

	ctrl := gomock.NewController(t)
	roleRep := getRoleRepMock(ctrl)
	defer ctrl.Finish()

	t.Run("Delete user by id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Create(ctx, &user)
		require.NoError(t, err)

		deletedUser, err := repos.Delete(ctx, user.UserID)
		require.NoError(t, err)
		assert.NotEqual(t, &user, deletedUser)
	})
	t.Run("Delete user by id with invalid id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Delete(ctx, user.UserID)
		require.Error(t, err)
	})
}

func Test_usersRepository_Update(t *testing.T) {
	t.Parallel()

	user := prepareUser()

	ctrl := gomock.NewController(t)
	roleRep := getRoleRepMock(ctrl)
	defer ctrl.Finish()

	t.Run("Update user by id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)

		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Create(ctx, &user)
		require.NoError(t, err)

		newUser := user
		newUser.Login = "different login"

		updatedUser, err := repos.Update(ctx, &newUser)
		require.NoError(t, err)
		assert.NotEqual(t, &user.Login, updatedUser.Login)
	})
	t.Run("Update user by id with unknown id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Update(ctx, &user)
		require.Error(t, err)
	})
}

func Test_usersRepository_Get(t *testing.T) {
	t.Parallel()

	user := prepareUser()

	ctrl := gomock.NewController(t)
	roleRep := getRoleRepMock(ctrl)
	defer ctrl.Finish()

	t.Run("Get user by id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &usersRepository{pg.DB(), roleRep, nil}

		_, err := repos.Create(ctx, &user)
		require.NoError(t, err)

		result, err := repos.GetByID(ctx, user.UserID)
		require.NoError(t, err)

		assert.Equal(t, &user, result)
	})

	t.Run("Get unknown user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &usersRepository{pg.DB(), roleRep, nil}

		result, err := repos.GetByID(ctx, user.UserID)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}
