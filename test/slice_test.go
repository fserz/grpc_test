package grpc_test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/v2/trees/redblacktree"
	"math"
	"slices"
	"testing"
)

func TestName(t *testing.T) {
	//query := [][]int{{2, 2}, {1, 2}, {3, 6}}
	query := []int{3, 6, 9, 12}
	for _, i := range query {
		fmt.Println(i)
	}
	//omitPasswordDemo()
	//intAndStringDemo()
	decoderDemo()
}

func decoderDemo() {
	// map[string]interface{} -> json string
	var m = make(map[string]interface{}, 1)
	m["count"] = 1 // int
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
	}
	fmt.Printf("str:%#v\n", string(b))
	// json string -> map[string]interface{}
	var m2 map[string]interface{}
	// 使用decoder方式反序列化，指定使用number类型
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(&m2)
	if err != nil {
		fmt.Printf("unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("value:%v\n", m2["count"]) // 1
	fmt.Printf("type:%T\n", m2["count"])  // json.Number
	// 将m2["count"]转为json.Number之后调用Int64()方法获得int64类型的值
	count, err := m2["count"].(json.Number).Int64()
	if err != nil {
		fmt.Printf("parse to int64 failed, err:%v\n", err)
		return
	}
	fmt.Printf("type:%T\n", int(count)) // int
}

type Card struct {
	ID    int64   `json:"id,string"`    // 添加string tag
	Score float64 `json:"score,string"` // 添加string tag
}

func intAndStringDemo() {
	jsonStr1 := `{"id": "1234567","score": "88.50"}`
	var c1 Card
	if err := json.Unmarshal([]byte(jsonStr1), &c1); err != nil {
		fmt.Printf("json.Unmarsha jsonStr1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("c1:%#v\n", c1) // c1:main.Card{ID:1234567, Score:88.5}
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type PublicUser struct {
	*User              // 匿名嵌套
	Password *struct{} `json:"password,omitempty"`
}

func omitPasswordDemo() {
	u1 := User{
		Name:     "七米",
		Password: "123456",
	}
	b, err := json.Marshal(PublicUser{User: &u1})
	if err != nil {
		fmt.Printf("json.Marshal u1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"name":"七米"}
}

func closestRoom(rooms [][]int, queries [][]int) []int {
	// 按照 size 从大到小排序
	slices.SortFunc(rooms, func(a, b []int) int { return b[1] - a[1] })

	q := len(queries)
	queryIds := make([]int, q)
	for i := range queryIds {
		queryIds[i] = i
	}
	// 按照 minSize 从大到小排序
	slices.SortFunc(queryIds, func(i, j int) int { return queries[j][1] - queries[i][1] })

	ans := make([]int, q)
	for i := range ans {
		ans[i] = -1
	}
	roomIds := redblacktree.New[int, struct{}]() // github.com/emirpasic/gods/v2/trees/redblacktree
	j := 0
	for _, i := range queryIds {
		preferredId, minSize := queries[i][0], queries[i][1]
		for j < len(rooms) && rooms[j][1] >= minSize {
			roomIds.Put(rooms[j][0], struct{}{})
			j++
		}

		diff := math.MaxInt
		// 左边的差
		if node, ok := roomIds.Floor(preferredId); ok {
			diff = preferredId - node.Key
			ans[i] = node.Key
		}
		// 右边的差
		if node, ok := roomIds.Ceiling(preferredId); ok && node.Key-preferredId < diff {
			ans[i] = node.Key
		}
	}
	return ans
}
