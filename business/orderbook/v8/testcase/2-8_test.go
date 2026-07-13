package testcase

import (
	"testing"
)

// cat stock_market.log | grep "2421-" |grep -a "】," | grep -av "【TradeHandler】"
func Test_2_8_A_1(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(90, "test_files/2-8/2-8-A-1.txt")
}

func Test_2_8_A_2(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(99, "test_files/2-8/2-8-A-2.txt")
}

func Test_2_8_C_1(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2403, "test_files/1-18-3-A.txt")
}

func Test_2_8_C_2(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2404, "test_files/1-18-4-A.txt")
}
