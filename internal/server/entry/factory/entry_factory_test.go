package factory

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/anoriar/gophkeeper/internal/server/entry/dto/request/sync"
	"github.com/anoriar/gophkeeper/internal/server/entry/entity"
	"github.com/anoriar/gophkeeper/internal/server/entry/enum"
	"github.com/anoriar/gophkeeper/internal/server/shared/services/uuid/mock_uuid_generator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEntryFactory_CreateEntryFromRequestItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGeneratorMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)

	type args struct {
		id          string
		requestItem sync.SyncRequestItem
		userID      string
	}

	decodedData, err := base64.StdEncoding.DecodeString("dL8e3WcogDHLFMwrCSPk9nZs8qXnWwBUupHiLuMuPaWDAuxBmUM/cH+Sv41fBb9OEf/AHjx0nx2yl5xewZM=")
	if err != nil {
		require.NoError(t, err)
	}

	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          entity.Entry
		err           error
	}{
		{
			name:          "success",
			mockBehaviour: func() {},
			args: args{
				id: "ab161651-ec2c-4cbb-a6c6-f8cf414e503d",
				requestItem: sync.SyncRequestItem{
					OriginalId: "65ad590d-77d6-49d4-a6e7-963d7b6f50a7",
					EntryType:  enum.Login,
					UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					Data:       "dL8e3WcogDHLFMwrCSPk9nZs8qXnWwBUupHiLuMuPaWDAuxBmUM/cH+Sv41fBb9OEf/AHjx0nx2yl5xewZM=",
					Meta:       json.RawMessage(""),
					IsDeleted:  false,
				},
				userID: "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
			},
			want: entity.Entry{
				Id:         "ab161651-ec2c-4cbb-a6c6-f8cf414e503d",
				OriginalId: "65ad590d-77d6-49d4-a6e7-963d7b6f50a7",
				UserId:     "b632eb93-0c31-4d6c-8fb9-282f3fb7e54e",
				EntryType:  enum.Login,
				UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
				Data:       decodedData,
				Meta:       json.RawMessage(""),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			f := &EntryFactory{
				uuidGen: uuidGeneratorMock,
			}
			got, err := f.CreateEntryFromRequestItem(tt.args.id, tt.args.requestItem, tt.args.userID)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("CreateEntryFromRequestItem() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			isEqual, err := got.Equals(tt.want)
			if err != nil {
				t.Errorf("CreateEntryFromRequestItem() check equal error %v", err)
			}
			if !isEqual {
				t.Errorf("CreateEntryFromRequestItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryFactory_CreateNewEntryFromRequestItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uuidGeneratorMock := mock_uuid_generator.NewMockUUIDGeneratorInterface(ctrl)

	type args struct {
		requestItem sync.SyncRequestItem
		userID      string
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          entity.Entry
		err           error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			f := &EntryFactory{
				uuidGen: uuidGeneratorMock,
			}
			got, err := f.CreateNewEntryFromRequestItem(tt.args.requestItem, tt.args.userID)
			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("CreateNewEntryFromRequestItem() error expectation: got = %v, want %v", err, tt.err)
				}
			}
			isEqual, err := got.Equals(tt.want)
			if err != nil {
				t.Errorf("CreateNewEntryFromRequestItem() check equal error %v", err)
			}
			if !isEqual {
				t.Errorf("CreateNewEntryFromRequestItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}