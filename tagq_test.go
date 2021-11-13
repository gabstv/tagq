package tagq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestQuerier struct {
	A        string              `json:"a"`
	Age      int                 `json:"age"`
	MapItems map[string]string   `header:"map_items" xml:"xmap-items"`
	Children []*TestQuerierChild `header:"hchildren" xml:"children"`
}

type TestQuerierChild struct {
	Marco      time.Time  `param:"marco"`
	Polo       *time.Time `param:"polo"`
	Scoreboard []int      `json:"scoreboard"`
}

func TestQ(t *testing.T) {
	n0 := time.Now().Add(time.Hour * -1)
	qx := &TestQuerier{
		A:   "real nice",
		Age: 42,
		MapItems: map[string]string{
			"waldo": "weldo",
			"fred":  "flintspears",
		},
		Children: []*TestQuerierChild{
			{
				Marco:      time.Now(),
				Polo:       &n0,
				Scoreboard: []int{1, 2, 3, 4, 5},
			},
			{
				Marco:      time.Now(),
				Polo:       &n0,
				Scoreboard: []int{6, 7, 8, 9, 10, 11},
			},
		},
	}
	assert.Equal(t, qx.A, Q(qx, "a").Str())
	assert.Equal(t, qx.Age, Q(qx, "age").Int())
	assert.Equal(t, qx.MapItems["waldo"], Q(qx, "xmap-items", "waldo").Str())
	assert.Equal(t, qx.MapItems["waldo"], Q(qx, "MapItems", "waldo").Str())
	assert.Equal(t, qx.MapItems["fred"], Q(qx, "map_items", "fred").Str())
	assert.Equal(t, qx.Children[0].Scoreboard[1], Q(qx, "hchildren", "0", "Scoreboard", "1").Int())
}
