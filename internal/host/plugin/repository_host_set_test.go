package plugin

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/boundary/internal/db"
	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/host/plugin/store"
	"github.com/hashicorp/boundary/internal/iam"
	"github.com/hashicorp/boundary/internal/kms"
	"github.com/hashicorp/boundary/internal/oplog"
	"github.com/hashicorp/boundary/internal/plugin/host"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestRepository_CreateSet(t *testing.T) {
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	kms := kms.TestKms(t, conn, wrapper)
	iamRepo := iam.TestRepo(t, conn, wrapper)
	_, prj := iam.TestScopes(t, iamRepo)
	plg := host.TestPlugin(t, conn, "create", "create")
	catalog := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 1)[0]

	tests := []struct {
		name      string
		in        *HostSet
		opts      []Option
		want      *HostSet
		wantIsErr errors.Code
	}{
		{
			name:      "nil-HostSet",
			wantIsErr: errors.InvalidParameter,
		},
		{
			name:      "nil-embedded-HostSet",
			in:        &HostSet{},
			wantIsErr: errors.InvalidParameter,
		},
		{
			name: "invalid-no-catalog-id",
			in: &HostSet{
				HostSet: &store.HostSet{},
			},
			wantIsErr: errors.InvalidParameter,
		},
		{
			name: "invalid-public-id-set",
			in: &HostSet{
				HostSet: &store.HostSet{
					CatalogId: catalog.PublicId,
					PublicId:  "abcd_OOOOOOOOOO",
				},
			},
			wantIsErr: errors.InvalidParameter,
		},
		{
			name: "valid-no-options",
			in: &HostSet{
				HostSet: &store.HostSet{
					CatalogId: catalog.PublicId,
				},
			},
			want: &HostSet{
				HostSet: &store.HostSet{
					CatalogId: catalog.PublicId,
				},
			},
		},
		{
			name: "valid-with-name",
			in: &HostSet{
				HostSet: &store.HostSet{
					CatalogId: catalog.PublicId,
					Name:      "test-name-repo",
				},
			},
			want: &HostSet{
				HostSet: &store.HostSet{
					CatalogId: catalog.PublicId,
					Name:      "test-name-repo",
				},
			},
		},
		{
			name: "valid-with-description",
			in: &HostSet{
				HostSet: &store.HostSet{
					CatalogId:   catalog.PublicId,
					Description: ("test-description-repo"),
				},
			},
			want: &HostSet{
				HostSet: &store.HostSet{
					CatalogId:   catalog.PublicId,
					Description: ("test-description-repo"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			repo, err := NewRepository(rw, rw, kms)
			require.NoError(err)
			require.NotNil(repo)
			got, err := repo.CreateSet(context.Background(), prj.GetPublicId(), tt.in, tt.opts...)
			if tt.wantIsErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantIsErr), err), "want err: %q got: %q", tt.wantIsErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)
			assert.Empty(tt.in.PublicId)
			require.NotNil(got)
			assert.True(strings.HasPrefix(got.GetPublicId(), HostSetPrefix))
			assert.NotSame(tt.in, got)
			assert.Equal(tt.want.Name, got.Name)
			assert.Equal(tt.want.Description, got.Description)
			assert.Equal(got.CreateTime, got.UpdateTime)
			assert.NoError(db.TestVerifyOplog(t, rw, got.PublicId, db.WithOperation(oplog.OpType_OP_TYPE_CREATE), db.WithCreateNotBefore(10*time.Second)))
		})
	}

	t.Run("invalid-duplicate-names", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		repo, err := NewRepository(rw, rw, kms)
		require.NoError(err)
		require.NotNil(repo)

		_, prj := iam.TestScopes(t, iamRepo)
		catalog := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 1)[0]

		in := &HostSet{
			HostSet: &store.HostSet{
				CatalogId: catalog.PublicId,
				Name:      "test-name-repo",
			},
		}

		got, err := repo.CreateSet(context.Background(), prj.GetPublicId(), in)
		require.NoError(err)
		require.NotNil(got)
		assert.True(strings.HasPrefix(got.GetPublicId(), HostSetPrefix))
		assert.NotSame(in, got)
		assert.Equal(in.Name, got.Name)
		assert.Equal(in.Description, got.Description)
		assert.Equal(got.CreateTime, got.UpdateTime)

		got2, err := repo.CreateSet(context.Background(), prj.GetPublicId(), in)
		assert.Truef(errors.Match(errors.T(errors.NotUnique), err), "want err code: %v got err: %v", errors.NotUnique, err)
		assert.Nil(got2)
	})

	t.Run("valid-duplicate-names-diff-catalogs", func(t *testing.T) {
		assert, require := assert.New(t), require.New(t)
		repo, err := NewRepository(rw, rw, kms)
		require.NoError(err)
		require.NotNil(repo)

		_, prj := iam.TestScopes(t, iamRepo)
		catalogs := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 2)

		catalogA, catalogB := catalogs[0], catalogs[1]

		in := &HostSet{
			HostSet: &store.HostSet{
				Name: "test-name-repo",
			},
		}
		in2 := in.clone()

		in.CatalogId = catalogA.PublicId
		got, err := repo.CreateSet(context.Background(), prj.GetPublicId(), in)
		require.NoError(err)
		require.NotNil(got)
		assert.True(strings.HasPrefix(got.GetPublicId(), HostSetPrefix))
		assert.NotSame(in, got)
		assert.Equal(in.Name, got.Name)
		assert.Equal(in.Description, got.Description)
		assert.Equal(got.CreateTime, got.UpdateTime)

		in2.CatalogId = catalogB.PublicId
		got2, err := repo.CreateSet(context.Background(), prj.GetPublicId(), in2)
		require.NoError(err)
		require.NotNil(got2)
		assert.True(strings.HasPrefix(got.GetPublicId(), HostSetPrefix))
		assert.NotSame(in2, got2)
		assert.Equal(in2.Name, got2.Name)
		assert.Equal(in2.Description, got2.Description)
		assert.Equal(got2.CreateTime, got2.UpdateTime)
	})
}

func TestRepository_LookupSet(t *testing.T) {
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	kms := kms.TestKms(t, conn, wrapper)
	iamRepo := iam.TestRepo(t, conn, wrapper)
	_, prj := iam.TestScopes(t, iamRepo)
	plg := host.TestPlugin(t, conn, "lookup", "lookup")

	catalog := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 1)[0]
	hostSet := TestSets(t, conn, catalog.PublicId, 1)[0]
	hostSetId, err := newHostSetId()
	require.NoError(t, err)

	tests := []struct {
		name      string
		in        string
		want      *HostSet
		wantIsErr errors.Code
	}{
		{
			name:      "with-no-public-id",
			wantIsErr: errors.InvalidParameter,
		},
		{
			name: "with-non-existing-host-set-id",
			in:   hostSetId,
		},
		{
			name:      "with-existing-host-set-id",
			in:        hostSet.PublicId,
			want:      hostSet,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			repo, err := NewRepository(rw, rw, kms)
			assert.NoError(err)
			require.NotNil(repo)
			got, err := repo.LookupSet(context.Background(), tt.in)
			if tt.wantIsErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantIsErr), err), "want err: %q got: %q", tt.wantIsErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)
			assert.EqualValues(tt.want, got)
		})
	}
}

func TestRepository_ListSets(t *testing.T) {
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	kms := kms.TestKms(t, conn, wrapper)
	iamRepo := iam.TestRepo(t, conn, wrapper)

	_, prj := iam.TestScopes(t, iamRepo)
	plg := host.TestPlugin(t, conn, "list", "list")
	catalogs := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 2)
	catalogA, catalogB := catalogs[0], catalogs[1]

	hostSets := TestSets(t, conn, catalogA.PublicId, 3)

	tests := []struct {
		name      string
		in        string
		opts      []Option
		want      []*HostSet
		wantIsErr errors.Code
	}{
		{
			name:      "with-no-catalog-id",
			wantIsErr: errors.InvalidParameter,
		},
		{
			name: "Catalog-with-no-host-sets",
			in:   catalogB.PublicId,
			want: []*HostSet{},
		},
		{
			name: "Catalog-with-host-sets",
			in:   catalogA.PublicId,
			want: hostSets,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			repo, err := NewRepository(rw, rw, kms)
			assert.NoError(err)
			require.NotNil(repo)
			got, err := repo.ListSets(context.Background(), tt.in, tt.opts...)
			if tt.wantIsErr != 0 {
				assert.Truef(errors.Match(errors.T(tt.wantIsErr), err), "want err: %q got: %q", tt.wantIsErr, err)
				assert.Nil(got)
				return
			}
			require.NoError(err)
			opts := []cmp.Option{
				cmpopts.SortSlices(func(x, y *HostSet) bool { return x.PublicId < y.PublicId }),
				protocmp.Transform(),
			}
			assert.Empty(cmp.Diff(tt.want, got, opts...))
		})
	}
}

func TestRepository_ListSets_Limits(t *testing.T) {
	conn, _ := db.TestSetup(t, "postgres")
	rw := db.New(conn)
	wrapper := db.TestWrapper(t)
	kms := kms.TestKms(t, conn, wrapper)
	iamRepo := iam.TestRepo(t, conn, wrapper)

	_, prj := iam.TestScopes(t, iamRepo)
	plg := host.TestPlugin(t, conn, "listlimit", "listlimit")
	catalog := TestCatalogs(t, conn, plg.GetPublicId(), prj.PublicId, 1)[0]
	count := 10
	hostSets := TestSets(t, conn, catalog.PublicId, count)

	tests := []struct {
		name     string
		repoOpts []Option
		listOpts []Option
		wantLen  int
	}{
		{
			name:    "With no limits",
			wantLen: count,
		},
		{
			name:     "With repo limit",
			repoOpts: []Option{WithLimit(3)},
			wantLen:  3,
		},
		{
			name:     "With negative repo limit",
			repoOpts: []Option{WithLimit(-1)},
			wantLen:  count,
		},
		{
			name:     "With List limit",
			listOpts: []Option{WithLimit(3)},
			wantLen:  3,
		},
		{
			name:     "With negative List limit",
			listOpts: []Option{WithLimit(-1)},
			wantLen:  count,
		},
		{
			name:     "With repo smaller than list limit",
			repoOpts: []Option{WithLimit(2)},
			listOpts: []Option{WithLimit(6)},
			wantLen:  6,
		},
		{
			name:     "With repo larger than list limit",
			repoOpts: []Option{WithLimit(6)},
			listOpts: []Option{WithLimit(2)},
			wantLen:  2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			repo, err := NewRepository(rw, rw, kms, tt.repoOpts...)
			assert.NoError(err)
			require.NotNil(repo)
			got, err := repo.ListSets(context.Background(), hostSets[0].CatalogId, tt.listOpts...)
			require.NoError(err)
			assert.Len(got, tt.wantLen)
		})
	}
}
