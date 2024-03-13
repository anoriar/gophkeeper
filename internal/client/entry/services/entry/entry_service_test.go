package entry

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	"github.com/anoriar/gophkeeper/internal/client/entry/factory/mock_entry_factory"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry/mock_entry_repository"
	"github.com/anoriar/gophkeeper/internal/client/entry/repository/entry_ext/mock_entry_ext_repository"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/encoder/mock_data_encryptor"
	"github.com/anoriar/gophkeeper/internal/client/shared/app/logger"
	sharedErrors "github.com/anoriar/gophkeeper/internal/client/shared/errors"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret"
	"github.com/anoriar/gophkeeper/internal/client/user/repository/secret/mock_secret_repository"
)

func TestLoginEntryService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	addEntryCommand := command.AddEntryCommand{
		EntryType: enum.Login,
		Data: dto.LoginData{
			Login:    "user",
			Password: "password",
		},
		Meta: []byte(""),
	}

	type args struct {
		ctx     context.Context
		command command.AddEntryCommand
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func(ctx context.Context, entryCommand command.AddEntryCommand)
		wantErr       error
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "cd06a579-311d-498e-aa01-d6ab589bf8bb",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				encryptedData := []byte("encrypted data")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromAddCmd(entryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(encryptedData, nil)
				entryMock.Data = encryptedData
				entryRepositoryMock.EXPECT().Add(ctx, entryMock).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "master pass not found error",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", secret.ErrMasterPasswordNotFound)
			},
			wantErr: secret.ErrMasterPasswordNotFound,
		},
		{
			name: "get master pass internal error",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "create entry internal error",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("12345", nil)
				entryFactoryMock.EXPECT().CreateFromAddCmd(addEntryCommand).Return(entity.Entry{}, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "encrypt internal error",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "cd06a579-311d-498e-aa01-d6ab589bf8bb",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}

				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromAddCmd(addEntryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(nil, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "save internal error",
			args: args{
				ctx:     context.Background(),
				command: addEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.AddEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "cd06a579-311d-498e-aa01-d6ab589bf8bb",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				encryptedData := []byte("encrypted data")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromAddCmd(entryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(encryptedData, nil)
				entryMock.Data = encryptedData
				entryRepositoryMock.EXPECT().Add(ctx, entryMock).Return(sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args.ctx, tt.args.command)
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			err := l.Add(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Add() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestLoginEntryService_Edit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	editEntryCommand := command.EditEntryCommand{
		Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
		EntryType: enum.Login,
		Data: dto.LoginData{
			Login:    "user",
			Password: "password",
		},
		Meta: []byte(""),
	}
	type args struct {
		ctx     context.Context
		command command.EditEntryCommand
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func(ctx context.Context, entryCommand command.EditEntryCommand)
		wantErr       error
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				encryptedData := []byte("encrypted data")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromEditCmd(entryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(encryptedData, nil)
				entryMock.Data = encryptedData
				entryRepositoryMock.EXPECT().Edit(ctx, entryMock).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "master pass not found error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", secret.ErrMasterPasswordNotFound)
			},
			wantErr: secret.ErrMasterPasswordNotFound,
		},
		{
			name: "get master pass internal error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "create entry internal error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("12345", nil)
				entryFactoryMock.EXPECT().CreateFromEditCmd(editEntryCommand).Return(entity.Entry{}, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "encrypt internal error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}

				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromEditCmd(editEntryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(nil, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "save internal error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				encryptedData := []byte("encrypted data")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromEditCmd(entryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(encryptedData, nil)
				entryMock.Data = encryptedData
				entryRepositoryMock.EXPECT().Edit(ctx, entryMock).Return(sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "save entry not found error",
			args: args{
				ctx:     context.Background(),
				command: editEntryCommand,
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.EditEntryCommand) {
				masterPass := "12345"
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				encryptedData := []byte("encrypted data")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryFactoryMock.EXPECT().CreateFromEditCmd(entryCommand).Return(entryMock, nil)
				encryptorMock.EXPECT().Encrypt(dataInBytes, masterPass).Return(encryptedData, nil)
				entryMock.Data = encryptedData
				entryRepositoryMock.EXPECT().Edit(ctx, entryMock).Return(sharedErrors.ErrEntryNotFound)
			},
			wantErr: sharedErrors.ErrEntryNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args.ctx, tt.args.command)
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			err := l.Edit(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Edit() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestLoginEntryService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		command command.DeleteEntryCommand
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		wantErr       error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			err := l.Delete(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestLoginEntryService_Detail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		command command.DetailEntryCommand
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          command_response.DetailEntryCommandResponse
		wantErr       error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			got, err := l.Detail(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Detail() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoginEntryService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		want          []command_response.ListEntryCommandResponse
		wantErr       error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			got, err := l.List(tt.args.ctx)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("List() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoginEntryService_Sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	entryFactoryMock := mock_entry_factory.NewMockEntryFactoryInterface(ctrl)
	entryRepositoryMock := mock_entry_repository.NewMockEntryRepositoryInterface(ctrl)
	secretRepositoryMock := mock_secret_repository.NewMockSecretRepositoryInterface(ctrl)
	encryptorMock := mock_data_encryptor.NewMockDataEncryptorInterface(ctrl)
	extRepositoryMock := mock_entry_ext_repository.NewMockEntryExtRepositoryInterface(ctrl)
	loggerMock, err := logger.Initialize("info")
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		command command.SyncEntryCommand
	}
	tests := []struct {
		name          string
		args          args
		mockBehaviour func()
		wantErr       error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()
			l := NewEntryService(
				entryFactoryMock,
				entryRepositoryMock,
				secretRepositoryMock,
				encryptorMock,
				extRepositoryMock,
				loggerMock,
			)
			err := l.Sync(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Sync() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}
