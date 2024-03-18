package entry

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"

	"github.com/stretchr/testify/assert"

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

func TestEntryService_Add(t *testing.T) {
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
		want          command_response.DetailEntryResponse
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
			want: command_response.DetailEntryResponse{
				Id:        "cd06a579-311d-498e-aa01-d6ab589bf8bb",
				EntryType: enum.Login,
				UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
				IsDeleted: false,
				Data: &dto.LoginData{
					Login:    "test",
					Password: "pass",
				},
				Meta: []byte(""),
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
			got, err := l.Add(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Add() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryService_Edit(t *testing.T) {
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
		want          command_response.DetailEntryResponse
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
			want: command_response.DetailEntryResponse{
				Id:        "ef77aba6-7ed4-421d-926a-93804ab96733",
				EntryType: enum.Login,
				UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
				IsDeleted: false,
				Data: &dto.LoginData{
					Login:    "test",
					Password: "pass",
				},
				Meta: []byte(""),
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
				dataInBytes := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
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
			got, err := l.Edit(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Edit() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Edit() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryService_Delete(t *testing.T) {
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
		mockBehaviour func(ctx context.Context, command command.DeleteEntryCommand)
		wantErr       error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				command: command.DeleteEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DeleteEntryCommand) {
				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entryMock, nil)
				entryRepositoryMock.EXPECT().Edit(ctx, entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: true,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "entry not found error",
			args: args{
				ctx: context.Background(),
				command: command.DeleteEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DeleteEntryCommand) {

				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entity.Entry{}, sharedErrors.ErrEntryNotFound)
			},
			wantErr: sharedErrors.ErrEntryNotFound,
		},
		{
			name: "get by id internal error",
			args: args{
				ctx: context.Background(),
				command: command.DeleteEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DeleteEntryCommand) {

				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entity.Entry{}, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "edit internal error",
			args: args{
				ctx: context.Background(),
				command: command.DeleteEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DeleteEntryCommand) {

				dataInBytes := []byte("test data")
				entryMock := entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entryMock, nil)
				entryRepositoryMock.EXPECT().Edit(ctx, entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: true,
					Data:      dataInBytes,
					Meta:      []byte(""),
				}).Return(sharedErrors.ErrInternalError)
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
			err := l.Delete(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Delete() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestEntryService_Detail(t *testing.T) {
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
		mockBehaviour func(ctx context.Context, command command.DetailEntryCommand)
		want          command_response.DetailEntryResponse
		wantErr       error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DetailEntryCommand) {
				masterPass := "1234"
				encryptedData := []byte("test data")
				entryMock := entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      encryptedData,
					Meta:      []byte(""),
				}
				decryptedData := []byte("{\"login\": \"test\", \"password\": \"pass\"}")
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entryMock, nil)
				encryptorMock.EXPECT().Decrypt(encryptedData, masterPass).Return(decryptedData, nil)
			},
			want: command_response.DetailEntryResponse{
				Id:        "225de857-71c5-452f-96f7-ff385d808083",
				EntryType: enum.Login,
				UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
				IsDeleted: false,
				Data: &dto.LoginData{
					Login:    "test",
					Password: "pass",
				},
				Meta: []byte(""),
			},
			wantErr: nil,
		},
		{
			name: "master pass not found error",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.DetailEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", secret.ErrMasterPasswordNotFound)
			},
			wantErr: secret.ErrMasterPasswordNotFound,
		},
		{
			name: "get master pass internal error",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.DetailEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("", sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "get by id not found error",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.DetailEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("12345", nil)
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entity.Entry{}, sharedErrors.ErrEntryNotFound)
			},
			wantErr: sharedErrors.ErrEntryNotFound,
		},
		{
			name: "get by id internal error",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, entryCommand command.DetailEntryCommand) {
				secretRepositoryMock.EXPECT().GetMasterPassword().Return("12345", nil)
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entity.Entry{}, sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "decrypt internal error",
			args: args{
				ctx: context.Background(),
				command: command.DetailEntryCommand{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
				},
			},
			mockBehaviour: func(ctx context.Context, command command.DetailEntryCommand) {
				masterPass := "1234"
				encryptedData := []byte("test data")
				entryMock := entity.Entry{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
					Data:      encryptedData,
					Meta:      []byte(""),
				}
				secretRepositoryMock.EXPECT().GetMasterPassword().Return(masterPass, nil)
				entryRepositoryMock.EXPECT().GetById(ctx, "225de857-71c5-452f-96f7-ff385d808083").Return(entryMock, nil)
				encryptorMock.EXPECT().Decrypt(encryptedData, masterPass).Return(nil, sharedErrors.ErrInternalError)
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
			got, err := l.Detail(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Detail() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Detail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryService_List(t *testing.T) {
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
	var tests = []struct {
		name          string
		args          args
		mockBehaviour func(ctx context.Context)
		want          []command_response.ListEntryCommandResponse
		wantErr       error
	}{
		{
			name: "success",
			args: args{ctx: context.Background()},
			mockBehaviour: func(ctx context.Context) {
				entryRepositoryMock.EXPECT().GetList(ctx).Return([]entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}, nil)
			},
			want: []command_response.ListEntryCommandResponse{
				{
					Id:        "225de857-71c5-452f-96f7-ff385d808083",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
				},
				{
					Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
					EntryType: enum.Login,
					UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
					IsDeleted: false,
				},
			},
		},
		{
			name: "get list internal error",
			args: args{ctx: context.Background()},
			mockBehaviour: func(ctx context.Context) {
				entryRepositoryMock.EXPECT().GetList(ctx).Return(nil, sharedErrors.ErrInternalError)
			},
			want:    nil,
			wantErr: sharedErrors.ErrInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args.ctx)
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
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntryService_Sync(t *testing.T) {
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
		mockBehaviour func(ctx context.Context, command command.SyncEntryCommand)
		wantErr       error
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				authToken := "cn8ewjf942tr49fehceo"
				secretRepositoryMock.EXPECT().GetAuthToken().Return(authToken, nil)
				entryRepositoryMock.EXPECT().GetList(ctx).Return([]entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}, nil)

				base64EncodedData1 := base64.StdEncoding.EncodeToString([]byte("data"))
				base64EncodedData2 := base64.StdEncoding.EncodeToString([]byte("data2"))

				syncResponse := entry_ext.SyncResponse{
					Items: []entry_ext.SyncResponseItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
					SyncType: enum.Login,
				}
				syncRequestMock := entry_ext.SyncRequest{
					SyncType: enum.Login,
					Items: []entry_ext.SyncRequestItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
				}
				extRepositoryMock.EXPECT().Sync(ctx, authToken, syncRequestMock).Return(syncResponse, nil)

				newEntries := []entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}
				entryFactoryMock.EXPECT().CreateFromSyncResponse(syncResponse).Return(newEntries, nil)
				entryRepositoryMock.EXPECT().Rewrite(ctx, newEntries).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "get auth token not found error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				secretRepositoryMock.EXPECT().GetAuthToken().Return("", secret.ErrTokenNotFound)
			},
			wantErr: secret.ErrTokenNotFound,
		},
		{
			name: "get auth token internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				secretRepositoryMock.EXPECT().GetAuthToken().Return("", sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "get auth token internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				secretRepositoryMock.EXPECT().GetAuthToken().Return("", sharedErrors.ErrInternalError)
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "get list internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				authToken := "f982hf8hwie"
				secretRepositoryMock.EXPECT().GetAuthToken().Return(authToken, nil)
				entryRepositoryMock.EXPECT().GetList(ctx).Return(nil, errors.New("error"))
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "sync internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				authToken := "f982hf8hwie"
				secretRepositoryMock.EXPECT().GetAuthToken().Return(authToken, nil)
				entryRepositoryMock.EXPECT().GetList(ctx).Return([]entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}, nil)
				base64EncodedData1 := base64.StdEncoding.EncodeToString([]byte("data"))
				base64EncodedData2 := base64.StdEncoding.EncodeToString([]byte("data2"))
				syncRequestMock := entry_ext.SyncRequest{
					SyncType: enum.Login,
					Items: []entry_ext.SyncRequestItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
				}
				extRepositoryMock.EXPECT().Sync(ctx, authToken, syncRequestMock).Return(entry_ext.SyncResponse{}, errors.New("error"))
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "create from sync response internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				authToken := "f982hf8hwie"
				secretRepositoryMock.EXPECT().GetAuthToken().Return(authToken, nil)
				entryRepositoryMock.EXPECT().GetList(ctx).Return([]entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}, nil)
				base64EncodedData1 := base64.StdEncoding.EncodeToString([]byte("data"))
				base64EncodedData2 := base64.StdEncoding.EncodeToString([]byte("data2"))
				syncResponse := entry_ext.SyncResponse{
					Items: []entry_ext.SyncResponseItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
					SyncType: enum.Login,
				}
				syncRequestMock := entry_ext.SyncRequest{
					SyncType: enum.Login,
					Items: []entry_ext.SyncRequestItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
				}
				extRepositoryMock.EXPECT().Sync(ctx, authToken, syncRequestMock).Return(syncResponse, nil)
				entryFactoryMock.EXPECT().CreateFromSyncResponse(syncResponse).Return(nil, errors.New("error"))
			},
			wantErr: sharedErrors.ErrInternalError,
		},
		{
			name: "rewrite internal error",
			args: args{
				ctx:     context.Background(),
				command: command.SyncEntryCommand{EntryType: enum.Login},
			},
			mockBehaviour: func(ctx context.Context, command command.SyncEntryCommand) {
				authToken := "f982hf8hwie"
				secretRepositoryMock.EXPECT().GetAuthToken().Return(authToken, nil)
				entryRepositoryMock.EXPECT().GetList(ctx).Return([]entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}, nil)
				base64EncodedData1 := base64.StdEncoding.EncodeToString([]byte("data"))
				base64EncodedData2 := base64.StdEncoding.EncodeToString([]byte("data2"))
				syncResponse := entry_ext.SyncResponse{
					Items: []entry_ext.SyncResponseItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
					SyncType: enum.Login,
				}
				syncRequestMock := entry_ext.SyncRequest{
					SyncType: enum.Login,
					Items: []entry_ext.SyncRequestItem{
						{
							OriginalId: "225de857-71c5-452f-96f7-ff385d808083",
							UpdatedAt:  time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData1,
							Meta:       []byte(""),
						},
						{
							OriginalId: "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
							UpdatedAt:  time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
							IsDeleted:  false,
							Data:       base64EncodedData2,
							Meta:       []byte(""),
						},
					},
				}
				extRepositoryMock.EXPECT().Sync(ctx, authToken, syncRequestMock).Return(syncResponse, nil)
				newEntries := []entity.Entry{
					{
						Id:        "225de857-71c5-452f-96f7-ff385d808083",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2023, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data"),
						Meta:      []byte(""),
					},
					{
						Id:        "60d016e5-eae1-49f6-bb00-7d4709a38f4c",
						EntryType: enum.Login,
						UpdatedAt: time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
						IsDeleted: false,
						Data:      []byte("data2"),
						Meta:      []byte(""),
					},
				}
				entryFactoryMock.EXPECT().CreateFromSyncResponse(syncResponse).Return(newEntries, nil)
				entryRepositoryMock.EXPECT().Rewrite(ctx, newEntries).Return(errors.New("error"))
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
			err := l.Sync(tt.args.ctx, tt.args.command)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Sync() error expectation: got = %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}
