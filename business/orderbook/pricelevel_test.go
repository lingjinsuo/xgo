package orderbook

import "testing"

func TestName(t *testing.T) {
	//t.Logf("%d", CalcIdx2Price(1, 11739))
	t.Logf("%d", CalcTickLevel(1, 9995000))
	t.Logf("%d", CalcTickLevel(1, 10000000))
}

func TestCalcTickLevel(t *testing.T) {
	idx := CalcTickLevel(2, 21400) - 1099
	t.Logf("%d", idx)
}
