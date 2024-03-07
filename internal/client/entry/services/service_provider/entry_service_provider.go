package service_provider

import (
	"context"

	"github.com/pkg/errors"

	"github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
	"github.com/anoriar/gophkeeper/internal/client/entry/enum"
	"github.com/anoriar/gophkeeper/internal/client/entry/services/entry"
)

type EntryServiceProvider struct {
	loginService entry.EntryServiceInterface
	cardService  entry.EntryServiceInterface
}

func NewEntryServiceProvider(loginService entry.EntryServiceInterface, cardService entry.EntryServiceInterface) *EntryServiceProvider {
	return &EntryServiceProvider{loginService: loginService, cardService: cardService}
}

func (sp *EntryServiceProvider) Add(ctx context.Context, cmd command.AddEntryCommand) error {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return err
	}
	err = service.Add(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (sp *EntryServiceProvider) Edit(ctx context.Context, cmd command.EditEntryCommand) error {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return err
	}
	err = service.Edit(ctx, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (sp *EntryServiceProvider) GetById(ctx context.Context, cmd command.DetailEntryCommand) (entity.Entry, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return entity.Entry{}, err
	}
	entryEntity, err := service.Detail(ctx, cmd)
	if err != nil {
		return entity.Entry{}, err
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

func (sp *EntryServiceProvider) GetList(ctx context.Context, cmd command.ListEntryCommand) ([]entity.Entry, error) {
	service, err := sp.getService(cmd.EntryType)
	if err != nil {
		return nil, err
	}
	entries, err := service.List(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (sp *EntryServiceProvider) getService(entryType enum.EntryType) (entry.EntryServiceInterface, error) {
	switch entryType {
	case enum.Login:
		return sp.loginService, nil
	case enum.Card:
		return sp.cardService, nil
	default:
		return nil, errors.New("not implemented cmd type")
	}
}