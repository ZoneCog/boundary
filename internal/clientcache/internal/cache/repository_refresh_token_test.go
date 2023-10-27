// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/boundary/api/authtokens"
	cachedb "github.com/hashicorp/boundary/internal/clientcache/internal/db"
	"github.com/hashicorp/boundary/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupRefreshToken(t *testing.T) {
	ctx := context.Background()
	s, err := cachedb.Open(ctx)
	require.NoError(t, err)

	r, err := NewRepository(ctx, s, &sync.Map{},
		mapBasedAuthTokenKeyringLookup(map[ringToken]*authtokens.AuthToken{}),
		sliceBasedAuthTokenBoundaryReader(nil))
	require.NoError(t, err)

	t.Run("nil user", func(t *testing.T) {
		_, err := r.lookupRefreshToken(ctx, nil, targetResourceType)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "user is nil")
	})

	t.Run("user id is empty", func(t *testing.T) {
		_, err := r.lookupRefreshToken(ctx, &user{Address: "addr"}, targetResourceType)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "user id is empty")
	})

	t.Run("resource type is invalid", func(t *testing.T) {
		_, err := r.lookupRefreshToken(ctx, &user{Id: "something", Address: "addr"}, resourceType("invalid"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "resource type is invalid")
	})

	t.Run("unknown user", func(t *testing.T) {
		got, err := r.lookupRefreshToken(ctx, &user{Id: "unkonwnUser", Address: "addr"}, targetResourceType)
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("no refresh token", func(t *testing.T) {
		known := &user{Id: "known", Address: "addr"}
		require.NoError(t, r.rw.Create(ctx, known))

		got, err := r.lookupRefreshToken(ctx, known, targetResourceType)
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("got refresh token", func(t *testing.T) {
		token := RefreshTokenValue("something")
		known := &user{Id: "withrefreshtoken", Address: "addr"}
		require.NoError(t, r.rw.Create(ctx, known))

		r.rw.DoTx(ctx, 1, db.ExpBackoff{}, func(r db.Reader, w db.Writer) error {
			require.NoError(t, upsertRefreshToken(ctx, w, known, targetResourceType, token))
			return nil
		})

		got, err := r.lookupRefreshToken(ctx, known, targetResourceType)
		assert.NoError(t, err)
		assert.Equal(t, token, got)
	})
}

func TestDeleteRefreshTokens(t *testing.T) {
	ctx := context.Background()
	s, err := cachedb.Open(ctx)
	require.NoError(t, err)

	r, err := NewRepository(ctx, s, &sync.Map{},
		mapBasedAuthTokenKeyringLookup(map[ringToken]*authtokens.AuthToken{}),
		sliceBasedAuthTokenBoundaryReader(nil))
	require.NoError(t, err)

	t.Run("nil user", func(t *testing.T) {
		err := r.deleteRefreshToken(ctx, nil, targetResourceType)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "user is nil")
	})

	t.Run("no user id", func(t *testing.T) {
		err := r.deleteRefreshToken(ctx, &user{Address: "addr"}, targetResourceType)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "user id is empty")
	})

	t.Run("invalid resource type", func(t *testing.T) {
		err := r.deleteRefreshToken(ctx, &user{Id: "id", Address: "addr"}, "this is invalid")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "resource type is invalid")
	})

	t.Run("success", func(t *testing.T) {
		u := &user{Id: "id", Address: "addr"}
		require.NoError(t, r.rw.Create(ctx, u))

		r.rw.DoTx(ctx, 1, db.ExpBackoff{}, func(r db.Reader, w db.Writer) error {
			require.NoError(t, upsertRefreshToken(ctx, w, u, targetResourceType, "token"))
			return nil
		})
		got, err := r.lookupRefreshToken(ctx, u, targetResourceType)
		require.NoError(t, err)
		require.NotEmpty(t, got)

		assert.NoError(t, r.deleteRefreshToken(ctx, u, targetResourceType))

		got, err = r.lookupRefreshToken(ctx, u, targetResourceType)
		require.NoError(t, err)
		require.Empty(t, got)
	})
}
