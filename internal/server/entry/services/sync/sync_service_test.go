package sync

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	syncRequestPkg "github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	syncResponsePkg "github.com/anoriar/gophkeeper/internal/server/entry/dto/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
	serverErrors "github.com/anoriar/gophkeeper/internal/server/entry/errors"
	"github.com/anoriar/gophkeeper/internal/server/entry/repository/entry_repository_mock"
	"github.com/anoriar/gophkeeper/internal/server/shared/app/db/mock"
	"github.com/anoriar/gophkeeper/internal/server/shared/app/logger"
	sharedErrors "github.com/anoriar/gophkeeper/internal/server/shared/errors"
	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid/mock_uuid_generator"
)

func TestSyncService_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)
	dbMock := mock.NewMockDatabaseInterface(ctrl)
	entryRepositoryMock := entry_repository_mock.NewMockEntryRepositoryInterface(ctrl)

	newItemData, err := base64.StdEncoding.DecodeString("L3lB71WXu7Jk25vSCsDmEKpsMYG6uqX+t8AyPZlkR1aaw7IhqEVoPaZ9Ds5vURD9fdqgzfRsEs3q6xUGwk4=")
	require.NoError(t, err)
	updateItemData, err := base64.StdEncoding.DecodeString("MX1mUs+puMP3FNlWITgzf5vS2JmcsVu/AivvxURLiQaPQJIxeVbF5/zGUBNVWuW5kzWhHKAi4E+gtoQ8Y9k=")
	require.NoError(t, err)
	deleteItemData, err := base64.StdEncoding.DecodeString("RKj38DKdE4z6qO0okC259mPyzENyiTd8UwQ7n3lIrpVLhmcgKkumi5fkygtxK8MaYv+Yy6wenhVIwDExfaU=")
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		request syncRequestPkg.SyncRequest
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          syncResponsePkg.SyncResponse
		err           error
	}{
		{
			name: "success add && update && delete items",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       newItemData,
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
							IsDeleted:  false,
						},
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       updateItemData,
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
							IsDeleted:  false,
						},
						{
							OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       deleteItemData,
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
							IsDeleted:  true,
						},
					},
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)

				updateItemOldData, err := base64.StdEncoding.DecodeString("ZGUj2QctZT3bkhxFH5zBkE/V6SK7ZOf1kzJL2+MHTU1aqXNWMLSheOTeot4hQZH4Yn0l5vbH0JsbbyO55rc=")
				require.NoError(t, err)

				firstGetUserEntriesMock := entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{Entries: []entity.Entry{
						{
							Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
							EntryType:  enum.Login,
							UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       updateItemOldData,
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
						},
						{
							Id:         "ffffc574-5eb0-4b3a-87af-93f2322f594e",
							EntryType:  enum.Login,
							UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
							OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       deleteItemData,
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
						},
					}}, nil)
				secondGetUserEntriesMock := entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{Entries: []entity.Entry{
						{
							Id:         "1675835b-f379-4121-a3f5-2b0abdb95c87",
							EntryType:  enum.Login,
							UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       newItemData,
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
						},
						{
							Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
							EntryType:  enum.Login,
							UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       updateItemData,
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
						},
					}}, nil)

				gomock.InOrder(
					firstGetUserEntriesMock,
					secondGetUserEntriesMock,
				)

				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				newItemUuidMock := uuidGenMock.EXPECT().NewString().Return("1675835b-f379-4121-a3f5-2b0abdb95c87")

				gomock.InOrder(
					newItemUuidMock,
				)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().AddEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "1675835b-f379-4121-a3f5-2b0abdb95c87",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
						UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						Data:       newItemData,
						Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
					},
				})
				entryRepositoryMock.EXPECT().UpdateEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       updateItemData,
						Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
					},
				})
				entryRepositoryMock.EXPECT().DeleteEntries(gomock.Any(), []string{
					"ffffc574-5eb0-4b3a-87af-93f2322f594e",
				})
				tx.EXPECT().Commit().Return(nil)
			},
			want: syncResponsePkg.SyncResponse{
				Items: []syncResponsePkg.SyncResponseItem{
					{
						OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
						UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						Data:       "L3lB71WXu7Jk25vSCsDmEKpsMYG6uqX+t8AyPZlkR1aaw7IhqEVoPaZ9Ds5vURD9fdqgzfRsEs3q6xUGwk4=",
						Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
					},
					{
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       "MX1mUs+puMP3FNlWITgzf5vS2JmcsVu/AivvxURLiQaPQJIxeVbF5/zGUBNVWuW5kzWhHKAi4E+gtoQ8Y9k=",
						Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
					},
				},
				SyncType: enum.Login,
			},
		},
		{
			name: "validation errors",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "",
							UpdatedAt:  time.Now(),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
					},
				},
			},
			mockBehaviour: func() {},
			want:          syncResponsePkg.SyncResponse{},
			err:           serverErrors.ErrSyncRequestNotValid,
		},
		{
			name: "get user entries first time internal error",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte("L3lB71WXu7Jk25vSCsDmEKpsMYG6uqX+t8AyPZlkR1aaw7IhqEVoPaZ9Ds5vURD9fdqgzfRsEs3q6xUGwk4="),
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
							IsDeleted:  false,
						},
					},
				},
			},
			mockBehaviour: func() {
				entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{}, errors.New("error"))
			},
			want: syncResponsePkg.SyncResponse{},
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "execute sync entries internal error",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte("L3lB71WXu7Jk25vSCsDmEKpsMYG6uqX+t8AyPZlkR1aaw7IhqEVoPaZ9Ds5vURD9fdqgzfRsEs3q6xUGwk4="),
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
							IsDeleted:  false,
						},
					},
				},
			},
			mockBehaviour: func() {
				entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{}, nil)
				uuidGenMock.EXPECT().NewString().Return("4c572e63-44d7-42dc-b87b-977d437732d4")
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(nil, errors.New("error"))
			},
			want: syncResponsePkg.SyncResponse{},
			err:  sharedErrors.ErrInternalError,
		},
		{
			name: "execute sync entries conflict error",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       newItemData,
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
							IsDeleted:  false,
						},
					},
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{}, nil)
				uuidGenMock.EXPECT().NewString().Return("1675835b-f379-4121-a3f5-2b0abdb95c87")
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().AddEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "1675835b-f379-4121-a3f5-2b0abdb95c87",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
						UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						Data:       newItemData,
						Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
					},
				}).Return(sharedErrors.ErrConflict)
			},
			want: syncResponsePkg.SyncResponse{},
			err:  sharedErrors.ErrConflict,
		},
		{
			name: "get user entries second time internal error",
			args: args{
				ctx: context.Background(),
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       newItemData,
							Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
							IsDeleted:  false,
						},
					},
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)

				getUserEntriesFirstMock := entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{}, nil)

				getUserEntriesSecondMock := entryRepositoryMock.EXPECT().GetEntriesByUserIDAndType(gomock.Any(), "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e", enum.Login).
					Return(collection.EntryCollection{}, errors.New("error"))

				gomock.InOrder(
					getUserEntriesFirstMock,
					getUserEntriesSecondMock,
				)

				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				newItemUuidMock := uuidGenMock.EXPECT().NewString().Return("1675835b-f379-4121-a3f5-2b0abdb95c87")

				gomock.InOrder(
					newItemUuidMock,
				)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().AddEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "1675835b-f379-4121-a3f5-2b0abdb95c87",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
						UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						Data:       newItemData,
						Meta:       []byte(`{"prop1": "valueProp1", "prop2": "valueProp2"}`),
					},
				})
				tx.EXPECT().Commit().Return(nil)
			},
			want: syncResponsePkg.SyncResponse{},
			err:  sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			s := NewSyncService(entryRepositoryMock, uuidGenMock, dbMock, loggerMock)
			got, err := s.Sync(tt.args.ctx, tt.args.request)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("Sync() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Sync() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_executeSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)
	dbMock := mock.NewMockDatabaseInterface(ctrl)
	entryRepositoryMock := entry_repository_mock.NewMockEntryRepositoryInterface(ctrl)

	type args struct {
		ctx            context.Context
		newEntries     []entity.Entry
		updatedEntries []entity.Entry
		deletedIds     []string
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				newEntries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				},
				updatedEntries: []entity.Entry{
					{
						Id:         "af66f6a8-f4f3-4759-991a-1e3d61b7b87d",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "814f0222-3df4-4217-89dd-0f0ea3f549b8",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				},
				deletedIds: []string{
					"1ecbfa8b-4697-4803-903c-856d40047bf6",
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().AddEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}).Return(nil)

				entryRepositoryMock.EXPECT().UpdateEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "af66f6a8-f4f3-4759-991a-1e3d61b7b87d",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "814f0222-3df4-4217-89dd-0f0ea3f549b8",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}).Return(nil)
				entryRepositoryMock.EXPECT().DeleteEntries(gomock.Any(), []string{
					"1ecbfa8b-4697-4803-903c-856d40047bf6",
				}).Return(nil)
				tx.EXPECT().Commit().Return(nil)
			},
			wantErr: false,
		},
		{
			name: "add entries error",
			args: args{
				ctx: context.Background(),
				newEntries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().AddEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "update entries error",
			args: args{
				ctx: context.Background(),
				updatedEntries: []entity.Entry{
					{
						Id:         "af66f6a8-f4f3-4759-991a-1e3d61b7b87d",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "814f0222-3df4-4217-89dd-0f0ea3f549b8",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().UpdateEntries(gomock.Any(), []entity.Entry{
					{
						Id:         "af66f6a8-f4f3-4759-991a-1e3d61b7b87d",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "814f0222-3df4-4217-89dd-0f0ea3f549b8",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "delete entries error",
			args: args{
				ctx:            context.Background(),
				newEntries:     []entity.Entry{},
				updatedEntries: []entity.Entry{},
				deletedIds:     []string{"nfs8dfjh234yfvc"},
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				entryRepositoryMock.EXPECT().DeleteEntries(gomock.Any(), []string{"nfs8dfjh234yfvc"}).Return(errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "commit error",
			args: args{
				ctx: context.Background(),
			},
			mockBehaviour: func() {
				tx := mock.NewMockDBTransactionInterface(ctrl)
				dbMock.EXPECT().BeginTransaction(gomock.Any()).Return(tx, nil)
				tx.EXPECT().Rollback().Return(nil)
				tx.EXPECT().Commit().Return(errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			s := NewSyncService(entryRepositoryMock, uuidGenMock, dbMock, loggerMock)
			if err := s.executeSync(tt.args.ctx, tt.args.newEntries, tt.args.updatedEntries, tt.args.deletedIds); (err != nil) != tt.wantErr {
				t.Errorf("executeSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncService_getDeletedIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)
	dbMock := mock.NewMockDatabaseInterface(ctrl)
	entryRepositoryMock := entry_repository_mock.NewMockEntryRepositoryInterface(ctrl)

	type args struct {
		request     syncRequestPkg.SyncRequest
		userEntries collection.EntryCollection
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "success",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
						{
							OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
					{
						Id:         "ffffc574-5eb0-4b3a-87af-93f2322f594e",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}},
			},
			want: []string{
				"ffffc574-5eb0-4b3a-87af-93f2322f594e",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSyncService(entryRepositoryMock, uuidGenMock, dbMock, loggerMock)
			if got := s.getDeletedIds(tt.args.request, tt.args.userEntries); !assert.Equal(t, got, tt.want) {
				t.Errorf("getDeletedIds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_getNewItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)
	dbMock := mock.NewMockDatabaseInterface(ctrl)
	entryRepositoryMock := entry_repository_mock.NewMockEntryRepositoryInterface(ctrl)

	type args struct {
		request     syncRequestPkg.SyncRequest
		userEntries collection.EntryCollection
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          []entity.Entry
	}{
		{
			name: "success add 2 new items",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
						{
							OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}},
			},
			mockBehaviour: func() {
				gomock.InOrder(
					uuidGenMock.EXPECT().NewString().Return("dfd0a1b2-ca49-4715-9900-504f06cc9e5b"),
					uuidGenMock.EXPECT().NewString().Return("ffffc574-5eb0-4b3a-87af-93f2322f594e"),
				)
			},
			want: []entity.Entry{
				{
					Id:         "dfd0a1b2-ca49-4715-9900-504f06cc9e5b",
					EntryType:  enum.Login,
					UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
					UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					Data:       []byte{},
					Meta:       []byte(""),
				},
				{
					Id:         "ffffc574-5eb0-4b3a-87af-93f2322f594e",
					EntryType:  enum.Login,
					UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
					UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
					Data:       []byte{},
					Meta:       []byte(""),
				},
			},
		},
		{
			name: "success try add isDeleted item",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{}},
			},
			mockBehaviour: func() {},
			want:          []entity.Entry{},
		},
	}
	for _, tt := range tests {
		tt.mockBehaviour()
		t.Run(tt.name, func(t *testing.T) {
			s := NewSyncService(entryRepositoryMock, uuidGenMock, dbMock, loggerMock)
			if got := s.getNewItems(tt.args.request, tt.args.userEntries); !assert.Equal(t, got, tt.want) {
				t.Errorf("getNewItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_getUpdatedItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGenMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)
	dbMock := mock.NewMockDatabaseInterface(ctrl)
	entryRepositoryMock := entry_repository_mock.NewMockEntryRepositoryInterface(ctrl)

	type args struct {
		request     syncRequestPkg.SyncRequest
		userEntries collection.EntryCollection
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          []entity.Entry
	}{
		{
			name: "success edit 1 item from 1 to add 1 to edit 1 to delete",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "0bc6c22e-d8b5-4057-9d28-eb5b7a233364",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  false,
						},
						{
							OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
							UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
					{
						Id:         "6b80f318-fac0-4ba7-889c-cbc2e2fc6223",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3453c579-9db6-4089-8ca3-1635a9887316",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}},
			},
			mockBehaviour: func() {},
			want: []entity.Entry{
				{
					Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
					EntryType:  enum.Login,
					UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
					UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					Data:       []byte{},
					Meta:       []byte(""),
				},
			},
		},
		{
			name: "success try to edit isDeleted item",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}},
			},
			mockBehaviour: func() {},
			want:          []entity.Entry{},
		},
		{
			name: "success try to edit item with older date",
			args: args{
				request: syncRequestPkg.SyncRequest{
					SyncType: enum.Login,
					UserID:   "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					Items: []syncRequestPkg.SyncRequestItem{
						{
							OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
							UpdatedAt:  time.Date(2021, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte{},
							Meta:       []byte(""),
							IsDeleted:  true,
						},
					},
				},
				userEntries: collection.EntryCollection{Entries: []entity.Entry{
					{
						Id:         "05dfdb32-3674-4381-be02-091e5e17080c",
						EntryType:  enum.Login,
						UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
						OriginalId: "3da6111c-6316-4993-aeff-74a2c3f345f9",
						UpdatedAt:  time.Date(2023, time.February, 10, 12, 0, 0, 0, time.UTC),
						Data:       []byte{},
						Meta:       []byte(""),
					},
				}},
			},
			mockBehaviour: func() {},
			want:          []entity.Entry{},
		},
	}
	for _, tt := range tests {
		tt.mockBehaviour()
		t.Run(tt.name, func(t *testing.T) {
			s := NewSyncService(entryRepositoryMock, uuidGenMock, dbMock, loggerMock)
			if got := s.getUpdatedItems(tt.args.request, tt.args.userEntries); !assert.Equal(t, got, tt.want) {
				t.Errorf("getUpdatedItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
