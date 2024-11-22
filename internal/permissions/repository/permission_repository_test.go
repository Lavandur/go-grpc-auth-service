package repository

import (
	"auth-service/testingdb"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	logger = logrus.New()
)

func Test_permissionRepository_GetRolePermissions(t *testing.T) {
	t.Parallel()

	t.Run("Get role permissions", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		pg := testingdb.NewWithIsolatedDatabase(t)
		repository := NewPermissionRepository(pg.DB(), logger)
		id := "3422b448-2460-4fd2-9183-8000de6f8343"
		expected := []string{"CAN_READ", "CAN_WRITE", "CAN_SEE"}

		_, err := repository.SetRolePermissions(ctx, id, expected)
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
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Set role permissions",
			fields: fields{
				db:     testingdb.NewWithIsolatedDatabase(t).DB(),
				logger: nil,
			},
			args: args{
				ctx:         context.Background(),
				id:          "3422b448-2460-4fd2-9183-8000de6f8343",
				permissions: []string{"CAN_READ", "CAN_WRITE"},
			},
			want:    true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &permissionRepository{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			got, err := p.SetRolePermissions(tt.args.ctx, tt.args.id, tt.args.permissions)
			if !tt.wantErr(t, err, fmt.Sprintf("SetRolePermissions(%v, %v)", tt.args.id, tt.args.permissions)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SetRolePermissions(%v, %v)", tt.args.id, tt.args.permissions)
		})
	}
}
