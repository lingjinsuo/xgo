package testcase

import (
	"testing"
)

// cat stock_market.log | grep "2421-" |grep -a "】," | grep -av "【TradeHandler】"
func Test_5_5H_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(199, "test_files/5-5H/5-5H-A.txt")
}
