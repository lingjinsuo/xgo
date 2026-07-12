package testcase

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Printf("111")
}

func Test_1_18_A(t *testing.T) {
	// 传递相对于 testcase 目录的路径
	RunOBTestCase(2421, "test_files/1-18-1-A.txt")
}
