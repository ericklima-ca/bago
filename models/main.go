package models

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Center{},
		&Order{},
		&PurchaseOrder{},
		&Customer{},
		&Product{},
		&Sell{},
		&PurchaseOrderStatus{},
		&OrderStatus{},
	}
}
