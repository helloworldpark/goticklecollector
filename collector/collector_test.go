package collector

import (
	"testing"
)

func TestCollectorCoinone(t *testing.T) {
	coinoneCol := coinoneCollector{}
	coins := coinoneCol.collect()
	if len(coins) != 11 {
		t.Errorf("Expected 11 coins but got %d coins", len(coins))
	}
}

func TestCollectorGopax(t *testing.T) {
	gopaxcol := gopaxCollector{"btc"}
	coins := gopaxcol.collect()
	if len(coins) != 1 {
		t.Errorf("Expected 1 coins but got %d coins", len(coins))
	}
}
