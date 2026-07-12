package orderbook

import (
	"gitlab.chuangzhen-sh.net/golang/xgo/logging/applogger"
)

var PriceLevelsIndexByIdxV1 []uint32
var PriceLevelsIndexByIdxV2 []uint32

func init() {
	initV1()
	initV2()
}

func initV1() {
	PriceLevelsIndexByIdxV1 = make([]uint32, 0)

	price := uint32(10) // 0.01
	for price < 10000000 {
		PriceLevelsIndexByIdxV1 = append(PriceLevelsIndexByIdxV1, price)
		price += CalcTickSizeV1(price)
	}

	applogger.Debug("PriceLevelsIndexByIdxV1 len=%d", len(PriceLevelsIndexByIdxV1))
}

func initV2() {
	PriceLevelsIndexByIdxV2 = make([]uint32, 0)

	price := uint32(10) // 0.01
	for price < 10000000 {
		PriceLevelsIndexByIdxV2 = append(PriceLevelsIndexByIdxV2, price)
		price += CalcTickSizeV2(price)
	}

	applogger.Debug("PriceLevelsIndexByIdxV2 len=%d", len(PriceLevelsIndexByIdxV2))
}

func CalcIdx2Price(ver int, idx uint32) uint32 {
	if ver == 1 {
		return PriceLevelsIndexByIdxV1[idx]
	}
	return PriceLevelsIndexByIdxV2[idx]
}

func CalcTickLevel(ver int, price uint32) uint32 {
	if ver == 1 {
		return CalcTickLevelV1(price)
	}
	return CalcTickLevelV2(price)
}

func CalcIdx2PriceV1(idx uint32) uint32 {
	return PriceLevelsIndexByIdxV1[idx]
}

func CalcIdx2PriceV2(idx uint32) uint32 {
	return PriceLevelsIndexByIdxV2[idx]
}

func CalcTickSizeV1(price uint32) uint32 {
	if price < 250 {
		return 1
	} // 0.01~0.25 0.001
	if price < 500 {
		return 5
	} // 0.25~0.5 0.005
	if price < 10000 {
		return 10
	} // 0.5~10.00 0.01
	if price < 20000 {
		return 20
	} // 10.00~20.00 0.02
	if price < 100000 {
		return 50
	} // 20~100 0.05
	if price < 200000 {
		return 100
	} // 100~200 0.1
	if price < 500000 {
		return 200
	} // 200~500 0.2
	if price < 1000000 {
		return 500
	} // 500~1000 0.5
	if price < 2000000 {
		return 1000
	} // 1000~2000 1
	if price < 5000000 {
		return 2000
	} // 2000~5000 2
	if price < 9995000 {
		return 5000
	} // 5000~9995 5
	return 5000
}

func CalcTickLevelV1(price uint32) uint32 {
	if price <= 250 { // 0.250
		return (price - 10) / 1
	}
	if price <= 500 { // 0.500
		return ((price - 250) / 5) + 240
	}
	if price <= 10000 { // 10.000
		return ((price - 500) / 10) + 290
	}
	if price <= 20000 {
		return ((price - 10000) / 20) + 1240
	}
	if price <= 100000 {
		return ((price - 20000) / 50) + 1740
	}
	if price <= 200000 {
		return ((price - 100000) / 100) + 3340
	}
	if price <= 500000 {
		return ((price - 200000) / 200) + 4340
	}
	if price <= 1000000 {
		return ((price - 500000) / 500) + 5840
	}
	if price <= 2000000 {
		return ((price - 1000000) / 1000) + 6840
	}
	if price <= 5000000 {
		return ((price - 2000000) / 2000) + 7840
	}
	if price <= 9995000 {
		return ((price - 5000000) / 5000) + 9340 // +1500 (5000000-2000000)/2000=1500
	}
	return 0
}

// 20250804 港交所新版的tick size规则
func CalcTickSizeV2(price uint32) uint32 {
	if price < 250 {
		return 1
	} // 0.01~0.25 0.001
	if price < 500 {
		return 5
	} // 0.25~0.5 0.005
	if price < 10000 {
		return 10
	} // 0.5~10.00 0.01
	if price < 20000 {
		return 10
	} // 10.00~20.00 0.01
	if price < 50000 {
		return 20
	} // 10.00~20.00 0.02
	if price < 100000 {
		return 50
	} // 20~100 0.05
	if price < 200000 {
		return 100
	} // 100~200 0.1
	if price < 500000 {
		return 200
	} // 200~500 0.2
	if price < 1000000 {
		return 500
	} // 500~1000 0.5
	if price < 2000000 {
		return 1000
	} // 1000~2000 1
	if price < 5000000 {
		return 2000
	} // 2000~5000 2
	if price < 9995000 {
		return 5000
	} // 5000~9995 5
	return 5000
}

// 20250804 港交所新版的tick size规则
func CalcTickLevelV2(price uint32) uint32 {
	if price <= 250 { // 0.250
		return (price - 10) / 1
	}
	if price <= 500 { // 0.500
		return ((price - 250) / 5) + 240
	}
	if price <= 10000 { // 10.000
		return ((price - 500) / 10) + 290
	} // +50 (500-250)/5=50
	if price <= 20000 {
		return ((price - 10000) / 20) + 1240
	} // +950 (10000-500)/10=950
	if price <= 50000 {
		return ((price - 10000) / 20) + 1740
	} // +500 (20000-10000)/20=500
	if price <= 100000 {
		return ((price - 20000) / 50) + 3240
	} // +1500 (50000-20000)/20=1500
	if price <= 200000 {
		return ((price - 100000) / 100) + 4240
	} // +1000 (100000-50000)/50=1000
	if price <= 500000 {
		return ((price - 200000) / 200) + 5240
	} // +1000 (200000-100000)/100=1000
	if price <= 1000000 {
		return ((price - 500000) / 500) + 6740
	} // +1500 (500000-200000)/200=1500
	if price <= 2000000 {
		return ((price - 1000000) / 1000) + 7740
	} // +1000 (1000000-500000)/500=1000
	if price <= 5000000 {
		return ((price - 2000000) / 2000) + 8740
	} // +1000 (2000000-1000000)/1000=1000
	if price <= 9995000 {
		return ((price - 5000000) / 5000) + 10240
	} // +1500 (5000000-2000000)/2000=1500
	return 0
}
