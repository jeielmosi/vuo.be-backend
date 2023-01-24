package adapters

type PostAtHashWrapper struct {
	hash string
}

func NewPostAtHashWrapperWrapper(hash string) *PostAtHashWrapper {
	return &PostAtHashWrapper{
		hash,
	}
}

func (this *PostAtHashWrapper) work (repository *ShortenBulkRepository) 
(
	*RepositoryDTO[K],
	error,
) {
	err :=  repository.PostAtHash(this.hash)
	return nil, err
}
