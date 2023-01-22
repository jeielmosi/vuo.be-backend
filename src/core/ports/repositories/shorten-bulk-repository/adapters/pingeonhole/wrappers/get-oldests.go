package adapters

type GetOldestsWrapper struct {
	size uint
}

func NewGetOldestsWrapper(size uint) *GetOldestsWrapper {
	return &GetOldestsWrapper{
		size,
	}
}

func (this *GetOldestsWrapper) work(repository *ShortenBulkRepository) 
(
	map[string]*PigeonholeDTO[K], 
	error,
) {
	return repository.GetOldests(this.size)
}
