package service_provider

import (
	"context"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command_response"

	"github.com/pkg/errors"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/entry"
)

type EntryServiceProvider struct {
	loginService entry.EntryServiceInterface
	cardService  entry.EntryServiceInterface
	textService  entry.EntryServiceInterface
	binService   entry.EntryServiceInterface
}

func NewEntryServiceProvider(loginService entry.EntryServiceInterface, cardService entry.EntryServiceInterface, textService entry.EntryServiceInterface, binService entry.EntryServiceInterface) *EntryServiceProvider {
	return &EntryServiceProvider{
		loginService: loginService,
		cardService:  cardService,
		textService:  textService,
		binService:   binService,
	}
}

func (sp *EntryServiceProvider) Add(ctx context.Context, cmd command.AddEntryCommand) (command_response.DetailEntryResponse, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	responseEntry, err := service.Add(ctx, cmd)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	return responseEntry, nil
}

func (sp *EntryServiceProvider) Edit(ctx context.Context, cmd command.EditEntryCommand) (command_response.DetailEntryResponse, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	responseEntry, err := service.Edit(ctx, cmd)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	return responseEntry, nil
}

func (sp *EntryServiceProvider) Detail(ctx context.Context, cmd command.DetailEntryCommand) (command_response.DetailEntryResponse, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	entryEntity, err := service.Detail(ctx, cmd)
	if err != nil {
		return command_response.DetailEntryResponse{}, err
	}
	return entryEntity, nil
}

func (sp *EntryServiceProvider) Delete(ctx context.Context, cmd command.DeleteEntryCommand) error {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return err
	}
	err = service.Delete(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (sp *EntryServiceProvider) GetList(ctx context.Context, cmd command.ListEntryCommand) ([]command_response.ListEntryCommandResponse, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return nil, err
	}
	entries, err := service.List(ctx)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (sp *EntryServiceProvider) Sync(ctx context.Context, cmd command.SyncEntryCommand) error {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return err
	}
	err = service.Sync(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (sp *EntryServiceProvider) getService(entryType enum.EntryType) (entry.EntryServiceInterface, error) {
	switch entryType {
	case enum.Login:
		return sp.loginService, nil
	case enum.Card:
		return sp.cardService, nil
	case enum.Text:
		return sp.textService, nil
	case enum.Bin:
		return sp.binService, nil
	default:
		return nil, errors.New("not implemented cmd type")
	}
}
