package repository

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/testingdb"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	logger = logrus.New()
)

func Test_roleRepository_Create(t *testing.T) {
	t.Parallel()

	role := models.Role{
		RoleID:      "3422b448-2460-4fd2-9183-8000de6f8343",
		Name:        "some_name",
		Description: map[string]string{"ru": "ADMIN"},
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
	}

	t.Run("Create role", func(t *testing.T) {
		t.Parallel()

		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		create, err := repos.Create(context.Background(), &role)
		assert.NoError(t, err)

		assert.Equal(t, &role, create)
	})

	t.Run("Create role with same id", func(t *testing.T) {
		t.Parallel()

		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Create(context.Background(), &role)
		assert.NoError(t, err)

		_, err = repos.Create(context.Background(), &role)
		assert.Error(t, err)
	})
}

func Test_roleRepository_Update(t *testing.T) {
	t.Parallel()

	role := models.Role{
		RoleID:      "3422b448-2460-4fd2-9183-8000de6f8343",
		Name:        "some_name",
		Description: map[string]string{"ru": "ADMIN"},
	}

	t.Run("Update role", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Create(ctx, &role)
		assert.NoError(t, err)

		role.Name = "some name"
		updated, err := repos.Update(ctx, &role)
		assert.NoError(t, err)
		assert.NotEqual(t, &role.Name, updated.Name)
	})

	t.Run("Update role with unknown id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Update(ctx, &role)
		assert.Error(t, err)
	})
}

func Test_roleRepository_Delete(t *testing.T) {
	t.Parallel()

	role := models.Role{
		RoleID:      "3422b448-2460-4fd2-9183-8000de6f8343",
		Name:        "some_name",
		Description: map[string]string{"ru": "ADMIN"},
	}

	t.Run("Delete role", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Create(ctx, &role)
		assert.NoError(t, err)

		err = repos.Delete(ctx, role.RoleID)
		assert.NoError(t, err)
	})

	t.Run("Delete role with unknown id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		err := repos.Delete(ctx, "some unused id")
		assert.Error(t, err)
	})
}

func Test_roleRepository_Getting(t *testing.T) {
	t.Parallel()

	role := models.Role{
		RoleID:      "3422b448-2460-4fd2-9183-8000de6f8343",
		Name:        "some_name",
		Description: map[string]string{"ru": "ADMIN"},
		CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
	}

	t.Run("Get role by id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Create(ctx, &role)
		assert.NoError(t, err)

		got, err := repos.GetByID(ctx, role.RoleID)
		assert.NoError(t, err)
		assert.Equal(t, &role, got)
	})

	t.Run("Get role with unknown id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.GetByID(ctx, role.RoleID)
		assert.Error(t, err)
	})

	t.Run("Get role by name", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		_, err := repos.Create(ctx, &role)
		assert.NoError(t, err)

		got, err := repos.GetByName(ctx, role.Name)
		assert.NoError(t, err)
		assert.Equal(t, &role, got)
	})

	t.Run("Get a list of roles", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repos := &roleRepository{pg.DB(), logger}

		for i := 0; i < 10; i++ {
			role := &models.Role{
				RoleID:      uuid.New().String(),
				Name:        "some_name",
				Description: map[string]string{"ru": "ADMIN"},
				CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			}

			_, err := repos.Create(ctx, role)
			assert.NoError(t, err)
		}

		orderBy := "created_at"
		var offset uint = 5
		var limit uint = 10
		pagination := &common.Pagination{
			OrderBy: &orderBy,
			Offset:  &offset,
			Size:    &limit,
		}

		got, err := repos.GetList(ctx, nil, nil)
		assert.NoError(t, err)
		assert.Len(t, got, 10)

		got, err = repos.GetList(ctx, nil, pagination)
		assert.NoError(t, err)
		assert.Len(t, got, 5)
	})
}
