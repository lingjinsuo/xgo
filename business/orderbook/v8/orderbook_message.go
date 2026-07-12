package v8

// AddOrder 新增订单消息
type AddOrder struct {
	OrderId      uint64
	SecurityCode uint32
	Price        uint32
	Quantity     uint32
	Side         uint16
}

// ModifyOrder 修改订单消息
type ModifyOrder struct {
	OrderId  uint64
	Quantity uint32
	Side     uint16
}

// DeleteOrder 删除订单消息
type DeleteOrder struct {
	OrderId uint64
	Side    uint16
}
