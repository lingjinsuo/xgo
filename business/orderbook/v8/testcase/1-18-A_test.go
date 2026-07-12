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

func Test_1_18_3_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2403, "test_files/1-18-3-A.txt")
}

func Test_1_18_4_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2404, "test_files/1-18-4-A.txt")
}

func Test_1_18_5_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2404, "test_files/1-18-5-A.txt")
}

func Test_1_18_6_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(3414, "test_files/1-18-6-A.txt")
}

func Test_1_18_7_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2921, "test_files/1-18-7-A.txt")
}

func Test_1_18_8_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(5201, "test_files/1-18-8-A.txt")
}
