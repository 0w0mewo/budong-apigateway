package server

type Service interface {
	RequestSetu(num int, isR18 bool) error
	GetSetuFromDB(id int) ([]byte, error)
	GetInventory(page, pageLimit uint64) ([]*SetuInfo, error)
	RandomSetu() ([]byte, error)
	Count() uint64
	Shutdown()
}
