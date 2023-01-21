package adapters

type GetByHashWrapper struct {
	repository *ShortenBulkRepository
	hash       string
}

func NewGetHashWrapper(repository *ShortenBulkRepository, hash string) *GetByHashWrapper {
	return &GetByHashWrapper{
		repository,
		hash,
	}
}

func work(this *GetByHashWrapper) (PigeonholeDTO[K], error) {
	return this.repository.GetByHash(this.hash)
}
