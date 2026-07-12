package testcase

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Printf("111")
}

// cat stock_market.log | grep "2421-" |grep -a "】," | grep -av "【TradeHandler】"
func Test_1_18_1_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2421, "test_files/1-18-1-A.txt")
}

func Test_1_18_2_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(3205, "test_files/1-18-2-A.txt")
}
