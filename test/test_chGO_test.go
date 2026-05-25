package news

import (
	"encoding/json"
	gproto "google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func testStat() Stat {
	return Stat{
		ID:      42,
		UID:     42,
		Type:    42,
		Name:    "name_value",
		Source:  "source_value",
		IP:      "127.0.0.1",
		Image:   true,
		Stack:   []string{"one", "two", "three"},
		Uints:   []uint{1, 2, 3},
		Uints8:  []uint8{1, 2, 3},
		Uints16: []uint16{1, 2, 3},
		Uints32: []uint32{1, 2, 3},
		Uints64: []uint64{1, 2, 3},
		Ints8:   []int8{1, 2, 3},
		Ints16:  []int16{1, 2, 3},
		Ints32:  []int32{1, 2, 3},
		Ints64:  []int64{1, 2, 3},
	}
}

func testStatProto() StatProto {
	return StatProto{
		Id:      42,
		Uid:     42,
		Type:    42,
		Name:    "name_value",
		Source:  "source_value",
		Ip:      "127.0.0.1",
		Image:   true,
		Stack:   []string{"one", "two", "three"},
		Uints:   []uint64{1, 2, 3},
		Uints8:  []byte{1, 2, 3},
		Uints16: []uint32{1, 2, 3},
		Uints32: []uint32{1, 2, 3},
		Uints64: []uint64{1, 2, 3},
		Ints8:   []int32{1, 2, 3},
		Ints16:  []int32{1, 2, 3},
		Ints32:  []int32{1, 2, 3},
		Ints64:  []int64{1, 2, 3},
	}
}

type testStatStd Stat

func TestStatMarshal(t *testing.T) {
	item := testStat()
	b := item.Marshal()
	if len(b) == 0 {
		t.Fatal("empty marshal result")
	}
}

func TestStatUnmarshal(t *testing.T) {
	item := testStat()
	b := item.Marshal()
	res := StatUnmarshal(b)
	if res == nil {
		t.Fatal("unmarshal failed")
	}
	if !reflect.DeepEqual(item, *res) {
		t.Fatalf("unmarshal mismatch: %#v != %#v", item, *res)
	}
}

func TestStatProtoMarshal(t *testing.T) {
	item := testStatProto()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty proto marshal result")
	}
}

func TestStatProtoUnmarshal(t *testing.T) {
	item := testStatProto()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	var res StatProto
	err = gproto.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if !gproto.Equal(&item, &res) {
		t.Fatalf("proto unmarshal mismatch: %#v != %#v", item, res)
	}
}

func TestStatGotinyMarshal(t *testing.T) {
	item := testStat()
	b := item.MarshalGotiny()
	if len(b) == 0 {
		t.Fatal("empty gotiny marshal result")
	}
}

func TestStatGotinyUnmarshal(t *testing.T) {
	item := testStat()
	b := item.MarshalGotiny()
	res := UnmarshalStatGotiny(b)
	if !reflect.DeepEqual(item, res) {
		t.Fatalf("gotiny unmarshal mismatch: %#v != %#v", item, res)
	}
}

func BenchmarkStatMarshalStd(b *testing.B) {
	item := testStat()
	std := testStatStd(item)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(&std)
	}
}

func BenchmarkStatUnmarshalStd(b *testing.B) {
	item := testStat()
	std := testStatStd(item)
	data, err := json.Marshal(&std)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res testStatStd
		err = json.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStatMarshal(b *testing.B) {
	item := testStat()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.Marshal()
	}
}

func BenchmarkStatUnmarshal(b *testing.B) {
	item := testStat()
	data := item.Marshal()
	b.ReportAllocs()
	for b.Loop() {
		res := StatUnmarshal(data)
		if res == nil {
			b.Fatal("unmarshal failed")
		}
	}
}

func BenchmarkStatProtoMarshal(b *testing.B) {
	item := testStatProto()
	b.ReportAllocs()
	for b.Loop() {
		_, err := gproto.Marshal(&item)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStatProtoUnmarshal(b *testing.B) {
	item := testStatProto()
	data, err := gproto.Marshal(&item)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res StatProto
		err = gproto.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStatGotinyMarshal(b *testing.B) {
	item := testStat()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.MarshalGotiny()
	}
}

func BenchmarkStatGotinyUnmarshal(b *testing.B) {
	item := testStat()
	data := item.MarshalGotiny()
	b.ReportAllocs()
	for b.Loop() {
		_ = UnmarshalStatGotiny(data)
	}
}
