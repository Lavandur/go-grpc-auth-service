package repository

import (
	"auth-service/internal/common"
	"auth-service/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_permissionRepository_GetRolePermissions(t *testing.T) {
	t.Parallel()

	t.Run("Get role permissions", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testdb.NewWithIsolatedDatabase(t)
		repository := NewPermissionRepository(pg.DB())
		id := "3422b448-2460-4fd2-9183-8000de6f8343"
		expected := []string{"CAN_READ", "CAN_WRITE", "CAN_SEE"}

		err := repository.SetRolePermissions(ctx, id, expected)
		assert.NoError(t, err)

		permissions, err := repository.GetRolePermissions(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, permissions)
	})
}

func Test_permissionRepository_SetRolePermissions(t *testing.T) {
	type fields struct {
		db     *pgxpool.Pool
		logger *logrus.Logger
	}
	type args struct {
		ctx         context.Context
		id          string
		permissions []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Set role permissions",
			fields: fields{
				db:     testdb.NewWithIsolatedDatabase(t).DB(),
				logger: nil,
			},
			args: args{
				ctx:         context.Background(),
				id:          "3422b448-2460-4fd2-9183-8000de6f8343",
				permissions: []string{"CAN_READ", "CAN_WRITE"},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &permissionRepository{
				db: tt.fields.db,
			}
			err := p.SetRolePermissions(tt.args.ctx, tt.args.id, tt.args.permissions)
			if !tt.wantErr(t, err, fmt.Sprintf("SetRolePermissions(%v, %v)", tt.args.id, tt.args.permissions)) {
				return
			}
		})
	}
}

func Test_permissionRepository_OperationsWithPermissions(t *testing.T) {
	t.Parallel()

	t.Run("Add permissions", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testdb.NewWithIsolatedDatabase(t)
		repos := NewPermissionRepository(pg.DB())

		expected := getPermission()

		result, err := repos.AddPermission(ctx, expected)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Delete permission by ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testdb.NewWithIsolatedDatabase(t)
		repos := NewPermissionRepository(pg.DB())

		expected := getPermission()

		_, err := repos.AddPermission(ctx, expected)
		require.NoError(t, err)

		err = repos.DeletePermission(ctx, expected.PermissionID)
		require.NoError(t, err)
	})

	t.Run("Get permission by ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testdb.NewWithIsolatedDatabase(t)
		repos := NewPermissionRepository(pg.DB())

		expected := getPermission()

		_, err := repos.AddPermission(ctx, expected)
		require.NoError(t, err)

		got, err := repos.GetPermissionByID(ctx, expected.PermissionID)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Get permission-list", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testdb.NewWithIsolatedDatabase(t)
		repos := NewPermissionRepository(pg.DB())

		got, err := repos.GetPermissions(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, len(got))

		_, err = repos.AddPermission(ctx, getPermission())
		require.NoError(t, err)
		_, err = repos.AddPermission(ctx, getPermission())
		require.NoError(t, err)
		_, err = repos.AddPermission(ctx, getPermission())
		require.NoError(t, err)

		got, err = repos.GetPermissions(ctx)
		require.NoError(t, err)
		assert.Equal(t, 3, len(got))
	})
}

func getPermission() *models.Permission {
	return &models.Permission{
		PermissionID: uuid.New().String(),
		Title:        uuid.New().String(),
		Description: common.LocalizedString{
			"en": "CAN_READ_ROLES",
		},
	}
}
