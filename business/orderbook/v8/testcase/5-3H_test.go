package testcase

import (
	"testing"
)

// cat stock_market.log | grep "2421-" |grep -a "】," | grep -av "【TradeHandler】"
func Test_5_3H_1(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(149, "test_files/5-3H/5-3H-1.txt")
}
