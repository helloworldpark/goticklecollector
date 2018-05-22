package collector

import (
	"testing"
)

func TestCollectorCoinone(t *testing.T) {
	coinoneCol := CoinoneCollector{}
	coins := coinoneCol.Collect()
	if len(coins) != 11 {
		t.Errorf("Expected 11 coins but got %d coins", len(coins))
	} else {
		t.Errorf("SUCCESS %v", coins)
	}
}

func TestCollectorGopax(t *testing.T) {
	gopaxcol := GopaxCollector{"btc"}
	coins := gopaxcol.Collect()
	if len(coins) != 1 {
		t.Errorf("Expected 1 coins but got %d coins", len(coins))
	} else {
		t.Errorf("SUCCESS %v", coins)
	}
}
