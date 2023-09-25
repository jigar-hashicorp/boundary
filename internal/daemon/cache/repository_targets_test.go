// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"sync"
	"testing"

	"github.com/hashicorp/boundary/api/authtokens"
	"github.com/hashicorp/boundary/api/targets"
	"github.com/hashicorp/boundary/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

func TestRepository_refreshTargets(t *testing.T) {
	ctx := context.Background()
	s, err := Open(ctx)
	require.NoError(t, err)

	addr := "address"
	u := user{
		Id:      "u1",
		Address: addr,
	}
	at := &authtokens.AuthToken{
		Id:     "at_1",
		Token:  "at_1_token",
		UserId: u.Id,
	}
	kt := KeyringToken{KeyringType: "k", TokenName: "t", AuthTokenId: at.Id}
	atMap := map[ringToken]*authtokens.AuthToken{
		{"k", "t"}: at,
	}
	r, err := NewRepository(ctx, s, &sync.Map{}, mapBasedAuthTokenKeyringLookup(atMap), sliceBasedAuthTokenBoundaryReader(maps.Values(atMap)))
	require.NoError(t, err)
	require.NoError(t, r.AddKeyringToken(ctx, addr, kt))

	ts := []*targets.Target{
		{
			Id:                "ttcp_1",
			Name:              "name1",
			Type:              "tcp",
			SessionMaxSeconds: 111,
		},
		{
			Id:                "ttcp_2",
			Name:              "name2",
			Address:           "address2",
			Type:              "tcp",
			SessionMaxSeconds: 222,
		},
		{
			Id:                "ttcp_3",
			Name:              "name3",
			Address:           "address3",
			Type:              "tcp",
			SessionMaxSeconds: 333,
		},
	}
	cases := []struct {
		name          string
		u             *user
		targets       []*targets.Target
		wantCount     int
		errorContains string
	}{
		{
			name: "Success",
			u: &user{
				Id:      at.UserId,
				Address: addr,
			},
			targets:   ts,
			wantCount: len(ts),
		},
		{
			name: "repeated target with different values",
			u: &user{
				Address: addr,
				Id:      at.UserId,
			},
			targets: append(ts, &targets.Target{
				Id:   ts[0].Id,
				Name: "a different name",
			}),
			wantCount: len(ts),
		},
		{
			name:          "nil user",
			u:             nil,
			targets:       ts,
			errorContains: "user is nil",
		},
		{
			name: "missing user Id",
			u: &user{
				Address: addr,
			},
			targets:       ts,
			errorContains: "user id is missing",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := r.refreshTargets(ctx, tc.u, tc.targets)
			if tc.errorContains == "" {
				assert.NoError(t, err)
				rw := db.New(s.conn)
				var got []*Target
				require.NoError(t, rw.SearchWhere(ctx, &got, "true", nil))
				assert.Len(t, got, tc.wantCount)
			} else {
				assert.ErrorContains(t, err, tc.errorContains)
			}
		})
	}
}

func TestRepository_ListTargets(t *testing.T) {
	ctx := context.Background()
	s, err := Open(ctx)
	require.NoError(t, err)

	addr := "address"
	u1 := &user{
		Id:      "u1",
		Address: addr,
	}
	at1 := &authtokens.AuthToken{
		Id:     "at_1",
		Token:  "at_1_token",
		UserId: u1.Id,
	}
	kt1 := KeyringToken{KeyringType: "k1", TokenName: "t1", AuthTokenId: at1.Id}

	u2 := &user{
		Id:      "u2",
		Address: addr,
	}
	at2 := &authtokens.AuthToken{
		Id:     "at_2",
		Token:  "at_2_token",
		UserId: u2.Id,
	}
	kt2 := KeyringToken{KeyringType: "k2", TokenName: "t2", AuthTokenId: at2.Id}
	atMap := map[ringToken]*authtokens.AuthToken{
		{"k1", "t1"}: at1,
		{"k2", "t2"}: at2,
	}
	r, err := NewRepository(ctx, s, &sync.Map{}, mapBasedAuthTokenKeyringLookup(atMap), sliceBasedAuthTokenBoundaryReader(maps.Values(atMap)))
	require.NoError(t, err)
	require.NoError(t, r.AddKeyringToken(ctx, addr, kt1))
	require.NoError(t, r.AddKeyringToken(ctx, addr, kt2))

	t.Run("token is missing", func(t *testing.T) {
		l, err := r.ListTargets(ctx, "")
		assert.Nil(t, l)
		assert.ErrorContains(t, err, "auth token id is missing")
	})

	ts := []*targets.Target{
		{
			Id:                "ttcp_1",
			Name:              "name1",
			Address:           "address1",
			Type:              "tcp",
			SessionMaxSeconds: 111,
		},
		{
			Id:                "ttcp_2",
			Name:              "name2",
			Address:           "address2",
			Type:              "tcp",
			SessionMaxSeconds: 222,
		},
		{
			Id:                "ttcp_3",
			Name:              "name3",
			Address:           "address3",
			Type:              "tcp",
			SessionMaxSeconds: 333,
		},
	}
	require.NoError(t, r.refreshTargets(ctx, u1, ts))

	t.Run("wrong user gets no targets", func(t *testing.T) {
		l, err := r.ListTargets(ctx, kt2.AuthTokenId)
		assert.NoError(t, err)
		assert.Empty(t, l)
	})
	t.Run("correct token gets targets", func(t *testing.T) {
		l, err := r.ListTargets(ctx, kt1.AuthTokenId)
		assert.NoError(t, err)
		assert.Len(t, l, len(ts))
		assert.ElementsMatch(t, l, ts)
	})
}

func TestRepository_QueryTargets(t *testing.T) {
	ctx := context.Background()
	s, err := Open(ctx)
	require.NoError(t, err)

	addr := "address"
	u1 := &user{
		Id:      "u1",
		Address: addr,
	}
	at1 := &authtokens.AuthToken{
		Id:     "at_1",
		Token:  "at_1_token",
		UserId: u1.Id,
	}
	kt1 := KeyringToken{KeyringType: "k1", TokenName: "t1", AuthTokenId: at1.Id}

	u2 := &user{
		Id:      "u2",
		Address: addr,
	}
	at2 := &authtokens.AuthToken{
		Id:     "at_2",
		Token:  "at_2_token",
		UserId: u2.Id,
	}
	kt2 := KeyringToken{KeyringType: "k2", TokenName: "t2", AuthTokenId: at2.Id}
	atMap := map[ringToken]*authtokens.AuthToken{
		{"k1", "t1"}: at1,
		{"k2", "t2"}: at2,
	}
	r, err := NewRepository(ctx, s, &sync.Map{}, mapBasedAuthTokenKeyringLookup(atMap), sliceBasedAuthTokenBoundaryReader(maps.Values(atMap)))
	require.NoError(t, err)
	require.NoError(t, r.AddKeyringToken(ctx, addr, kt1))
	require.NoError(t, r.AddKeyringToken(ctx, addr, kt2))

	query := "name % name1 or name % name2"

	errorCases := []struct {
		name        string
		p           string
		query       string
		errContains string
	}{
		{
			name:        "auth token id is missing",
			p:           "",
			query:       query,
			errContains: "auth token id is missing",
		},
		{
			name:        "query is missing",
			p:           "authtokenid",
			errContains: "query is missing",
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			l, err := r.QueryTargets(ctx, tc.p, tc.query)
			assert.Nil(t, l)
			assert.ErrorContains(t, err, tc.errContains)
		})
	}

	ts := []*targets.Target{
		{
			Id:                "ttcp_1",
			Name:              "name1",
			Address:           "address1",
			Type:              "tcp",
			SessionMaxSeconds: 111,
		},
		{
			Id:                "ttcp_2",
			Name:              "name2",
			Address:           "address2",
			Type:              "tcp",
			SessionMaxSeconds: 222,
		},
		{
			Id:                "ttcp_3",
			Name:              "name3",
			Address:           "address3",
			Type:              "tcp",
			SessionMaxSeconds: 333,
		},
	}
	require.NoError(t, r.refreshTargets(ctx, u1, ts))

	t.Run("wrong token gets no targets", func(t *testing.T) {
		l, err := r.QueryTargets(ctx, kt2.AuthTokenId, query)
		assert.NoError(t, err)
		assert.Empty(t, l)
	})
	t.Run("correct token gets targets", func(t *testing.T) {
		l, err := r.QueryTargets(ctx, kt1.AuthTokenId, query)
		assert.NoError(t, err)
		assert.Len(t, l, 2)
		assert.ElementsMatch(t, l, ts[0:2])
	})
}
