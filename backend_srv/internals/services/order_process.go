package services

import "scm-simple-luke.com/dir/internals/db"

type OrderProcess struct {
	Store db.DBTX
}

func NewOrderProcess(store db.DBTX) *OrderProcess {
	return &OrderProcess{
		Store: store,
	}
}

// test
func (o *OrderProcess) ProcessData(key, newval string) error {
	// return o.Store.Update(key, newval, time.Now())
	return nil
}
