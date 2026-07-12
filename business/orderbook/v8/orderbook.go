package v8

import (
	"fmt"
	"sort"
)

// 买卖方向常量
const (
	Buy  uint16 = 0 // 买方向
	Sell uint16 = 1 // 卖方向
)

// 主动方向常量
const (
	ActiveDirectionNone uint16 = 0 // 未知或无法判断
	ActiveDirectionBuy  uint16 = 1 // 主动买（向上吃单）
	ActiveDirectionSell uint16 = 2 // 主动卖（向下吃单）
)

// 订单结构体
type Order struct {
	OrderID      uint64
	SecurityCode uint32
	Price        uint32 // 注意：Price是int32
	Quantity     uint32
	Side         uint16
}

// OrderBook 订单簿结构体
type OrderBook struct {
	SecurityCode uint32
	BuyOrders    map[uint64]*Order // OrderID -> Order (买盘)
	SellOrders   map[uint64]*Order // OrderID -> Order (卖盘)

	// 用于记录潜在的立即成交订单
	pendingImmediateOrders map[uint64]*Order // 可能立即成交的订单
	lastBid1BeforeCross    uint32            // 价格交叉前的买一价
	lastAsk1BeforeCross    uint32            // 价格交叉前的卖一价
}

// 价格水平汇总
type PriceLevel struct {
	Price    uint32
	Quantity uint32
}

// 创建新的订单簿
func New(securityCode uint32) *OrderBook {
	return &OrderBook{
		SecurityCode:           securityCode,
		BuyOrders:              make(map[uint64]*Order),
		SellOrders:             make(map[uint64]*Order),
		pendingImmediateOrders: make(map[uint64]*Order),
		lastBid1BeforeCross:    0,
		lastAsk1BeforeCross:    0,
	}
}

// Add 处理新增订单消息
func (ob *OrderBook) Add(msg *AddOrder) error {
	// 检查订单是否已存在
	if _, exists := ob.BuyOrders[msg.OrderId]; exists {
		return fmt.Errorf("order %d already exists in buy orders", msg.OrderId)
	}
	if _, exists := ob.SellOrders[msg.OrderId]; exists {
		return fmt.Errorf("order %d already exists in sell orders", msg.OrderId)
	}

	// 创建新订单
	order := &Order{
		OrderID:      msg.OrderId,
		SecurityCode: msg.SecurityCode,
		Price:        msg.Price,
		Quantity:     msg.Quantity,
		Side:         msg.Side,
	}

	// 获取当前买卖一档价格
	//bid1, _, ask1, _, hasBoth := ob.GetBestBidAsk()

	// 检查是否可能导致立即成交
	if msg.Side == Buy {
		// 买单：价格 >= 卖一价，可能导致立即成交
		/*if hasBoth && msg.Price >= ask1 {
			// 记录价格交叉前的买卖一价
			if ob.lastBid1BeforeCross == 0 {
				ob.lastBid1BeforeCross = bid1
				ob.lastAsk1BeforeCross = ask1
			}
			// 将订单标记为可能立即成交，但不加入买盘
			ob.pendingImmediateOrders[msg.OrderId] = order
			return nil
		}*/
		// 正常买单，加入买盘
		ob.BuyOrders[msg.OrderId] = order
	} else if msg.Side == Sell {
		// 卖单：价格 <= 买一价，可能导致立即成交
		/*if hasBoth && msg.Price <= bid1 {
			// 记录价格交叉前的买卖一价
			if ob.lastBid1BeforeCross == 0 {
				ob.lastBid1BeforeCross = bid1
				ob.lastAsk1BeforeCross = ask1
			}
			// 将订单标记为可能立即成交，但不加入卖盘
			ob.pendingImmediateOrders[msg.OrderId] = order
			return nil
		}*/
		// 正常卖单，加入卖盘
		ob.SellOrders[msg.OrderId] = order
	} else {
		return fmt.Errorf("unknown order side: %d", msg.Side)
	}

	return nil
}

// Modify 处理修改订单消息
func (ob *OrderBook) Modify(msg *ModifyOrder) error {
	// 首先检查是否在待成交订单中
	if order, exists := ob.pendingImmediateOrders[msg.OrderId]; exists {
		// 修改待成交订单的数量
		order.Quantity = msg.Quantity
		return nil
	}

	// 检查正常订单
	var order *Order
	var exists bool

	if msg.Side == Buy {
		order, exists = ob.BuyOrders[msg.OrderId]
	} else if msg.Side == Sell {
		order, exists = ob.SellOrders[msg.OrderId]
	} else {
		return fmt.Errorf("unknown order side: %d", msg.Side)
	}

	if !exists {
		return fmt.Errorf("order %d not found", msg.OrderId)
	}

	// 更新订单数量
	order.Quantity = msg.Quantity

	return nil
}

// Delete 处理删除订单消息
func (ob *OrderBook) Delete(msg *DeleteOrder) error {
	// 首先检查是否在待成交订单中
	if _, exists := ob.pendingImmediateOrders[msg.OrderId]; exists {
		delete(ob.pendingImmediateOrders, msg.OrderId)
		return nil
	}

	// 检查正常订单
	var exists bool

	if msg.Side == Buy {
		_, exists = ob.BuyOrders[msg.OrderId]
		if exists {
			delete(ob.BuyOrders, msg.OrderId)
			return nil
		}
	} else if msg.Side == Sell {
		_, exists = ob.SellOrders[msg.OrderId]
		if exists {
			delete(ob.SellOrders, msg.OrderId)
			return nil
		}
	} else {
		return fmt.Errorf("unknown order side: %d", msg.Side)
	}

	return fmt.Errorf("order %d not found", msg.OrderId)
}

// GetBestBidAsk 获取最优买价和卖价
func (ob *OrderBook) GetBestBidAsk() (uint32, uint32, uint32, uint32, bool) {
	buyLevels := ob.GetBuyPriceLevels()
	sellLevels := ob.GetSellPriceLevels()

	var bestBidPx, bestBidQty, bestAskPx, bestAskQty uint32
	hasBid := len(buyLevels) > 0
	hasAsk := len(sellLevels) > 0

	if hasBid {
		bestBidPx = buyLevels[0].Price
		bestBidQty = buyLevels[0].Quantity
	}
	if hasAsk {
		bestAskPx = sellLevels[0].Price
		bestAskQty = sellLevels[0].Quantity
	}

	return bestBidPx, bestBidQty, bestAskPx, bestAskQty, hasBid && hasAsk
}

// GetBuyPriceLevels 获取买盘价格水平汇总（按价格降序排列）
func (ob *OrderBook) GetBuyPriceLevels() []PriceLevel {
	return ob.getPriceLevels(ob.BuyOrders, true)
}

// GetSellPriceLevels 获取卖盘价格水平汇总（按价格升序排列）
func (ob *OrderBook) GetSellPriceLevels() []PriceLevel {
	return ob.getPriceLevels(ob.SellOrders, false)
}

// getPriceLevels 通用方法：获取价格水平汇总
func (ob *OrderBook) getPriceLevels(orders map[uint64]*Order, descending bool) []PriceLevel {
	// 按价格汇总数量
	priceMap := make(map[uint32]uint32)
	for _, order := range orders {
		priceMap[order.Price] += order.Quantity
	}

	// 转换为切片
	levels := make([]PriceLevel, 0, len(priceMap))
	for price, qty := range priceMap {
		levels = append(levels, PriceLevel{Price: price, Quantity: qty})
	}

	// 排序
	if descending {
		// 买盘：价格从高到低
		sort.Slice(levels, func(i, j int) bool {
			return levels[i].Price > levels[j].Price
		})
	} else {
		// 卖盘：价格从低到高
		sort.Slice(levels, func(i, j int) bool {
			return levels[i].Price < levels[j].Price
		})
	}

	return levels
}

// PrintOrderBook 打印订单簿
func (ob *OrderBook) PrintOrderBook(levels int) {
	fmt.Printf("\n=== 订单簿: %d ===\n", ob.SecurityCode)

	// 买盘
	buyLevels := ob.GetBuyPriceLevels()
	fmt.Println("买盘 (价格从高到低):")
	for i := 0; i < len(buyLevels) && i < levels; i++ {
		level := buyLevels[i]
		fmt.Printf("  %d \t %d\n", level.Price, level.Quantity)
	}

	// 卖盘
	sellLevels := ob.GetSellPriceLevels()
	fmt.Println("\n卖盘 (价格从低到高):")
	for i := 0; i < len(sellLevels) && i < levels; i++ {
		level := sellLevels[i]
		fmt.Printf("  %d \t %d\n", level.Price, level.Quantity)
	}

	// 最优报价
	bid, _, ask, _, hasBoth := ob.GetBestBidAsk()
	if hasBoth {
		fmt.Printf("\n最优报价: 买=%d, 卖=%d, 价差=%d\n", bid, ask, ask-bid)
	}

	// 待成交订单信息
	pendingCount := len(ob.pendingImmediateOrders)
	if pendingCount > 0 {
		fmt.Printf("待成交订单: %d 个\n", pendingCount)
		for _, order := range ob.pendingImmediateOrders {
			var sideStr string
			if order.Side == Buy {
				sideStr = "买"
			} else {
				sideStr = "卖"
			}
			fmt.Printf("  订单ID: %d, 方向: %s, 价格: %d, 数量: %d\n",
				order.OrderID, sideStr, order.Price, order.Quantity)
		}
	}

	fmt.Println("=====================")
}

// 获取订单数量统计
func (ob *OrderBook) GetOrderStats() (int, int, int) {
	return len(ob.BuyOrders), len(ob.SellOrders), len(ob.pendingImmediateOrders)
}

// GetPendingImmediateOrders 获取待成交订单列表（用于调试）
func (ob *OrderBook) GetPendingImmediateOrders() []*Order {
	orders := make([]*Order, 0, len(ob.pendingImmediateOrders))
	for _, order := range ob.pendingImmediateOrders {
		orders = append(orders, order)
	}
	return orders
}

// PrintOrderBookTable 打印订单簿表格（Bid-Ask格式）
func (ob *OrderBook) PrintOrderBookTable() {
	// 打印表头
	fmt.Printf("OrderID\tOrderType\tQuantity\tPrice\tPrice\tQuantity\tOrderType\tOrderID\n")

	// 打印买盘订单（按价格降序）
	buyOrders := ob.getSortedOrdersByPrice(ob.BuyOrders, true)
	for _, order := range buyOrders {
		price := float64(order.Price) / 1000.0
		fmt.Printf("%d\t2\t%d\t%.3f\n", order.OrderID, order.Quantity, price)
	}

	// 打印卖盘订单（按价格升序）
	sellOrders := ob.getSortedOrdersByPrice(ob.SellOrders, false)
	for _, order := range sellOrders {
		price := float64(order.Price) / 1000.0
		fmt.Printf("\t\t\t\t%.3f\t%d\t2\t%d\n", price, order.Quantity, order.OrderID)
	}
}

// getSortedOrdersByPrice 按价格排序订单（价格相同时按数量降序）
func (ob *OrderBook) getSortedOrdersByPrice(orders map[uint64]*Order, descending bool) []*Order {
	result := make([]*Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, order)
	}
	sort.Slice(result, func(i, j int) bool {
		if descending {
			if result[i].Price == result[j].Price {
				return result[i].Quantity > result[j].Quantity // 数量大的在前
			}
			return result[i].Price > result[j].Price
		}
		if result[i].Price == result[j].Price {
			return result[i].Quantity > result[j].Quantity // 数量大的在前
		}
		return result[i].Price < result[j].Price
	})
	return result
}
