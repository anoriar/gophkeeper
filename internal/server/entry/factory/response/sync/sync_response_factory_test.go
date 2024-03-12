package sync

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/anoriar/gophkeeper/internal/server/entry/dto/collection"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/response/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
)

func TestSyncResponseFactory_CreateSyncResponse(t *testing.T) {
	type args struct {
		entryCollection collection.EntryCollection
		syncType        enum.EntryType
	}
	tests := []struct {
		name string
		args args
		want sync.SyncResponse
	}{
		{
			name: "success",
			args: args{
				entryCollection: collection.EntryCollection{
					Entries: []entity.Entry{
						{
							Id:         "7f8b3f0d-321e-440e-9cb3-e9d80d6a9db2",
							OriginalId: "54493f7e-b64f-4831-8b38-691768a86d83",
							UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
							EntryType:  enum.Login,
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       []byte("0x2F7941EF5597BBB264DB9BD20AC0E610AA6C3181BABAA5FEB7C0323D996447569AC3B221A845683DA67D0ECE6F5110FD7DDAA0CDF46C12CDEAEB1506C24E"),
							Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
						},
					},
				},
				syncType: enum.Login,
			},
			want: sync.SyncResponse{
				Items: []sync.SyncResponseItem{
					{
						OriginalId: "54493f7e-b64f-4831-8b38-691768a86d83",
						UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						Data:       base64.StdEncoding.EncodeToString([]byte("0x2F7941EF5597BBB264DB9BD20AC0E610AA6C3181BABAA5FEB7C0323D996447569AC3B221A845683DA67D0ECE6F5110FD7DDAA0CDF46C12CDEAEB1506C24E")),
						Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
					},
				},
				SyncType: enum.Login,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SyncResponseFactory{}
			if got := f.CreateSyncResponse(tt.args.entryCollection, tt.args.syncType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSyncResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncResponseFactory_CreateSyncResponseItem(t *testing.T) {
	type args struct {
		entry entity.Entry
	}
	tests := []struct {
		name string
		args args
		want sync.SyncResponseItem
	}{
		{
			name: "success",
			args: args{
				entry: entity.Entry{
					Id:         "7f8b3f0d-321e-440e-9cb3-e9d80d6a9db2",
					OriginalId: "54493f7e-b64f-4831-8b38-691768a86d83",
					UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
					EntryType:  enum.Login,
					UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					Data:       []byte("0x2F7941EF5597BBB264DB9BD20AC0E610AA6C3181BABAA5FEB7C0323D996447569AC3B221A845683DA67D0ECE6F5110FD7DDAA0CDF46C12CDEAEB1506C24E"),
					Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
				},
			},
			want: sync.SyncResponseItem{
				OriginalId: "54493f7e-b64f-4831-8b38-691768a86d83",
				UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
				Data:       base64.StdEncoding.EncodeToString([]byte("0x2F7941EF5597BBB264DB9BD20AC0E610AA6C3181BABAA5FEB7C0323D996447569AC3B221A845683DA67D0ECE6F5110FD7DDAA0CDF46C12CDEAEB1506C24E")),
				Meta:       []byte(`{"key1": "value1", "key2": "value2"}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SyncResponseFactory{}
			if got := f.CreateSyncResponseItem(tt.args.entry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSyncResponseItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
