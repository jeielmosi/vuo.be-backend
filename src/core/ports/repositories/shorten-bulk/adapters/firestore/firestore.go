package firestore

import (
	"cloud.google.com/go/firestore"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkFirestore struct {
	envName string
}

// TODO: Verify if is OK
func (f *ShortenBulkFirestore) Get(hash string) (
	*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	client, err, ctx := getClient(f.envName)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	snapshot, err := client.Collection(ShortenBulkCollection).Doc(hash).Get(ctx)
	if err != nil || !snapshot.Exists() {
		return nil, err
	}

	return types.ToRepositoryDTO(snapshot.Data())
}

// TODO: Verify if is OK
func (f *ShortenBulkFirestore) GetOldest(size int) (
	map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity],
	error,
) {
	client, err, ctx := getClient(f.envName)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	iter := client.
		Collection(ShortenBulkCollection).
		OrderBy(types.CreatedAtField, firestore.Asc).
		Where(types.LockedField, "==", false).
		Limit(size).
		Documents(ctx)

	mp := map[string]*repositories.RepositoryDTO[entities.ShortenBulkEntity]{}
	for true {
		snapshot, err := iter.Next()
		if err != nil {
			break
		}
		dto, err := types.ToRepositoryDTO(snapshot.Data())
		if err != nil {
			break
		}

		mp[snapshot.Ref.ID] = dto
	}

	return mp, nil
}

func (f *ShortenBulkFirestore) Post(
	hash string,
	dto repositories.RepositoryDTO[entities.ShortenBulkEntity],
) error {
	client, err, ctx := getClient(f.envName)
	if err != nil {
		return err
	}
	defer client.Close()

	flatten := types.NewShortenBulkFlattenDTO(dto)
	_, err = client.Collection(ShortenBulkCollection).
		Doc(hash).
		Set(ctx, flatten)

	return err
}

// TODO: return a interface type
func NewShortenBulkFirestore(envName string) *ShortenBulkFirestore {
	return &ShortenBulkFirestore{
		envName,
	}
}

/*

TODO:

	IncrementClicks(hash string) error
	Lock(hash string) error
	Unlock(hash string) error
*/
