package firestore_shorten_bulk

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	entities "github.com/jei-el/vuo.be-backend/src/core/domain/shorten-bulk"
	shorten_bulk "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/interfaces"
	types "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/shorten-bulk/types"
	repositories "github.com/jei-el/vuo.be-backend/src/core/ports/repositories/types"
)

type ShortenBulkFirestore struct {
	envName string
}

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
			log.Println(err.Error())
			break
		}
		dto, err := types.ToRepositoryDTO(snapshot.Data())
		if err != nil {
			log.Println(err.Error())
			continue
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
	ref := client.Collection(ShortenBulkCollection).Doc(hash)
	fn := func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}

		lock, err := doc.DataAt(types.LockedField)
		currLock := lock.(bool)

		if flatten[types.LockedField].(bool) != currLock {
			return errors.New("Updating the lock status in a wrong way")
		}

		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	return client.RunTransaction(ctx, fn)
}

func (f *ShortenBulkFirestore) IncrementClicks(hash string) error {
	client, err, ctx := getClient(f.envName)
	if err != nil {

		return err
	}
	defer client.Close()

	ref := client.Collection(ShortenBulkCollection).Doc(hash)
	fn := func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}

		clicks, err := doc.DataAt(types.ClicksField)
		if err != nil {
			return err
		}

		flatten := types.ShortenBulkFlattenDTO{}
		flatten[types.ClicksField] = clicks.(int64) + 1

		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	return client.RunTransaction(ctx, fn)
}

func (f *ShortenBulkFirestore) updateLocked(hash string, locked bool) error {
	client, err, ctx := getClient(f.envName)
	if err != nil {
		return err
	}
	defer client.Close()

	ref := client.Collection(ShortenBulkCollection).Doc(hash)

	fn := func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}

		curr, err := doc.DataAt(types.LockedField)
		if err != nil {
			return err
		}

		if curr.(bool) == locked {
			log.Println(curr, locked)
			return errors.New("Document is already locked/unlocked")
		}

		flatten := types.ShortenBulkFlattenDTO{}
		flatten[types.LockedField] = locked

		return tx.Set(ref, flatten, firestore.MergeAll)
	}

	return client.RunTransaction(ctx, fn)
}

func (f *ShortenBulkFirestore) Lock(hash string) error {
	return f.updateLocked(hash, true)
}

func (f *ShortenBulkFirestore) Unlock(hash string) error {
	return f.updateLocked(hash, false)
}

func NewShortenBulkFirestore(envName string) shorten_bulk.ShortenBulkRepository {
	return &ShortenBulkFirestore{
		envName,
	}
}
