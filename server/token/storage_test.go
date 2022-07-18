package token_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"word_of_wisdom/server/token"
	"word_of_wisdom/test"
)

var (
	storageError = errors.New("error")
)

func TestNewOnetimeStorage_Get(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		ctx, ctxCancel = context.WithCancel(context.Background())
	)
	defer ctxCancel()

	tests := []struct {
		name string

		storageGet test.MockContract

		expected      uint
		expectedError error
	}{
		{
			name:       "ok",
			storageGet: test.MockContract{Param1: "key", Calls: 1},
			expected:   1,
		},
		{
			name:          "error, storage error",
			storageGet:    test.MockContract{Param1: "key", Value2: storageError, Calls: 1},
			expectedError: storageError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stor := NewMockStorage(ctrl)
			stor.EXPECT().
				Get(ctx, tc.storageGet.Param1).Return(tc.storageGet.Value1, tc.storageGet.Value2).Times(tc.storageGet.Calls)

			err := token.NewOnetimeStorage(stor).Get(ctx, "key")
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewOnetimeStorage_Put(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		ctx, ctxCancel = context.WithCancel(context.Background())
	)
	defer ctxCancel()

	tests := []struct {
		name string

		storagePut test.MockContract

		expectedError error
	}{
		{
			name:          "ok",
			storagePut:    test.MockContract{Param1: "key", Value1: nil, Calls: 1},
			expectedError: nil,
		},
		{
			name:          "err",
			storagePut:    test.MockContract{Param1: "key", Value1: storageError, Calls: 1},
			expectedError: storageError,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stor := NewMockStorage(ctrl)
			stor.EXPECT().
				Put(ctx, tc.storagePut.Param1, struct{}{}).Return(tc.storagePut.Value1).Times(tc.storagePut.Calls)

			err := token.NewOnetimeStorage(stor).Put(ctx, "key")
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestNewOnetimeStorage_Use(t *testing.T) {
	var (
		ctrl           = gomock.NewController(t)
		ctx, ctxCancel = context.WithCancel(context.Background())
	)
	defer ctxCancel()

	tests := []struct {
		name string

		storageGet    test.MockContract
		storageDelete test.MockContract

		expectedError error
	}{
		{
			name:          "ok",
			storageGet:    test.MockContract{Param1: "key", Value1: nil, Calls: 1},
			storageDelete: test.MockContract{Param1: "key", Value1: nil, Calls: 1},
			expectedError: nil,
		},
		{
			name:          "error, get not found",
			storageGet:    test.MockContract{Param1: "key", Value2: storageError, Calls: 1},
			expectedError: storageError,
		},
		{
			name:          "error, delete",
			storageGet:    test.MockContract{Param1: "key", Value1: nil, Calls: 1},
			storageDelete: test.MockContract{Param1: "key", Value1: errors.New("error"), Calls: 1},
			expectedError: errors.New("error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stor := NewMockStorage(ctrl)
			stor.EXPECT().
				Get(ctx, tc.storageGet.Param1).Return(tc.storageGet.Value1, tc.storageGet.Value2).Times(tc.storageGet.Calls)
			stor.EXPECT().
				Delete(ctx, tc.storageDelete.Param1).Return(tc.storageDelete.Value1).Times(tc.storageDelete.Calls)

			err := token.NewOnetimeStorage(stor).Use(ctx, "key")
			require.Equal(t, tc.expectedError, err)
		})
	}
}
