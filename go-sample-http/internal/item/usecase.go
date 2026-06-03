package item

type UseCase interface {
	List() ([]Item, error)
	Get(id int64) (*Item, error)
	Create(req CreateRequest) (*Item, error)
	Update(id int64, req UpdateRequest) (*Item, error)
	Delete(id int64) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) List() ([]Item, error) {
	return uc.repo.FindAll()
}

func (uc *useCase) Get(id int64) (*Item, error) {
	return uc.repo.FindByID(id)
}

func (uc *useCase) Create(req CreateRequest) (*Item, error) {
	return uc.repo.Create(req)
}

func (uc *useCase) Update(id int64, req UpdateRequest) (*Item, error) {
	return uc.repo.Update(id, req)
}

func (uc *useCase) Delete(id int64) error {
	return uc.repo.Delete(id)
}
