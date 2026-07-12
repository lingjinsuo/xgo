package v8

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	orderbook := New(700)

	// 测试买方向 (Buy=0)
	err := orderbook.Add(&AddOrder{
		SecurityCode: 700,
		OrderId:      1001,
		Price:        32000,
		Quantity:     100,
		Side:         0, // Buy
	})
	if err != nil {
		fmt.Printf("添加买单错误: %v\n", err)
	}

	// 测试卖方向 (Sell=1)
	err = orderbook.Add(&AddOrder{
		SecurityCode: 700,
		OrderId:      2001,
		Price:        32400,
		Quantity:     100,
		Side:         1, // Sell
	})
	if err != nil {
		fmt.Printf("添加卖单错误: %v\n", err)
	}

	// 获取最优报价
	bid, _, ask, _, hasBoth := orderbook.GetBestBidAsk()
	if hasBoth {
		fmt.Printf("最优报价: 买一价=%d, 卖一价=%d\n", bid, ask)
	}

	// 测试成交
	err = orderbook.Add(&AddOrder{
		SecurityCode: 700,
		OrderId:      1002,
		Price:        32400,
		Quantity:     100,
		Side:         0,
	})
	if err != nil {
		fmt.Printf("添加卖单错误: %v\n", err)
	}

	if err != nil {
		fmt.Printf("处理成交错误: %v\n", err)
	}

}
