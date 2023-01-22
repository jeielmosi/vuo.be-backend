package adapters

type GetByHashWrapper struct {
	hash string
}

func NewGetHashWrapper(hash string) *GetByHashWrapper {
	return &GetByHashWrapper{
		hash,
	}
}

func (this *GetByHashWrapper) work (repository *ShortenBulkRepository) 
(
	*PigeonholeDTO[K],
	error,
) {
	return repository.GetByHash(this.hash)
}
