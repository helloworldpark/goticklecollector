package gatherer

import (
	"fmt"
	"testing"

	"github.com/helloworldpark/goticklecollector/collector"
)

func TestGatherer(t *testing.T) {
	coinoneGW := GiveWork(collector.CoinoneCollector{}, 3)
	dfBundle, _ := Gather(coinoneGW)

	for _, b := range dfBundle {
		coin := <-b.Channel()
		fmt.Println(fmt.Sprintf("%v", coin))
	}
	t.Errorf("Finished")
}

func TestMerger(t *testing.T) {
	m1 := make([]CoinGateway, 10)
	m2 := make([]CoinGateway, 20)
	m3 := make([]CoinGateway, 30)

	mm := MergeGateways(m1, m2, m3)

	if len(mm) != 60 {
		t.Error("Failed")
	}
}
