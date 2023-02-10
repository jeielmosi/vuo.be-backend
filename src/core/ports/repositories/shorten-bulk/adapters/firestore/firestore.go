package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	repository_helpers "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/helpers"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

func getClient(envName string) (*firestore.Client, error, context.Context) {
	ctx := context.Background()

	app, err := getApp(envName)
	if err != nil {
		return nil, err, ctx
	}

	client, err := app.Firestore(ctx)
	return client, err, ctx
}

type ShortenBulkFirestore struct {
	envName string
}

// TODO: use const to standardize
func (f *ShortenBulkFirestore) Get(hash string) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	client, err, ctx := getClient(f.envName)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	snapshot, err := client.Collection(Collection).Doc(hash).Get(ctx)
	if err != nil || !snapshot.Exists() {
		return nil, err
	}

	data := snapshot.Data()

	createdAtStr := data[CreatedAtField].(string)
	createdAt, err := repository_helpers.NewTimeFrom10NanosecondsString(&createdAtStr)
	if err != nil {
		return nil, err
	}

	updatedAtStr := data[UpdatedAtField].(string)
	updatedAt, err := repository_helpers.NewTimeFrom10NanosecondsString(&updatedAtStr)
	if err != nil {
		return nil, err
	}

	dto := &repositories.RepositoryDTO[entities.ShortenBulkEntity]{
		Entity: entities.NewShortenBulkEntity(
			data[URLField].(string),
			data[ClicksField].(uint64),
		),
		CreatedAt: *createdAt,
		Locked:    data[LockedField].(bool),
		UpdatedAt: *updatedAt,
	}

	return dto, err
}

// TODO: return a interface type
func NewShortenBulkFirestore(envName string) *ShortenBulkFirestore {
	return &ShortenBulkFirestore{
		envName,
	}
}

/*

TODO:

	GetOldests(size uint) (map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity], error)
	Post(hash string, shorten_bulk repositories.RepositoryDTO[entities.ShortenBulkEntity]) error
	IncrementClicks(hash string) error
	Lock(hash string) error
	Unlock(hash string) error
*/
