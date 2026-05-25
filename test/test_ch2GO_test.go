package news

import (
	"encoding/json"
	gproto "google.golang.org/protobuf/proto"
	"reflect"
	"testing"
)

func testStat2() Stat2 {
	return Stat2{
		ID:     42,
		UID:    42,
		Type:   42,
		Name:   "name_value",
		Source: "source_value",
		IP:     "127.0.0.1",
		Image:  true,
	}
}

func testStat2Proto() Stat2Proto {
	return Stat2Proto{
		Id:     42,
		Uid:    42,
		Type:   42,
		Name:   "name_value",
		Source: "source_value",
		Ip:     "127.0.0.1",
		Image:  true,
	}
}

type testStat2Std Stat2

func TestStat2Marshal(t *testing.T) {
	item := testStat2()
	b := item.Marshal()
	if len(b) == 0 {
		t.Fatal("empty marshal result")
	}
}

func TestStat2Unmarshal(t *testing.T) {
	item := testStat2()
	b := item.Marshal()
	res := Stat2Unmarshal(b)
	if res == nil {
		t.Fatal("unmarshal failed")
	}
	if !reflect.DeepEqual(item, *res) {
		t.Fatalf("unmarshal mismatch: %#v != %#v", item, *res)
	}
}

func TestStat2ProtoMarshal(t *testing.T) {
	item := testStat2Proto()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty proto marshal result")
	}
}

func TestStat2ProtoUnmarshal(t *testing.T) {
	item := testStat2Proto()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	var res Stat2Proto
	err = gproto.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if !gproto.Equal(&item, &res) {
		t.Fatalf("proto unmarshal mismatch: %#v != %#v", item, res)
	}
}

func TestStat2GotinyMarshal(t *testing.T) {
	item := testStat2()
	b := item.MarshalGotiny()
	if len(b) == 0 {
		t.Fatal("empty gotiny marshal result")
	}
}

func TestStat2GotinyUnmarshal(t *testing.T) {
	item := testStat2()
	b := item.MarshalGotiny()
	res := UnmarshalStat2Gotiny(b)
	if !reflect.DeepEqual(item, res) {
		t.Fatalf("gotiny unmarshal mismatch: %#v != %#v", item, res)
	}
}

func BenchmarkStat2MarshalStd(b *testing.B) {
	item := testStat2()
	std := testStat2Std(item)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(&std)
	}
}

func BenchmarkStat2UnmarshalStd(b *testing.B) {
	item := testStat2()
	std := testStat2Std(item)
	data, err := json.Marshal(&std)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res testStat2Std
		err = json.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStat2Marshal(b *testing.B) {
	item := testStat2()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.Marshal()
	}
}

func BenchmarkStat2Unmarshal(b *testing.B) {
	item := testStat2()
	data := item.Marshal()
	b.ReportAllocs()
	for b.Loop() {
		res := Stat2Unmarshal(data)
		if res == nil {
			b.Fatal("unmarshal failed")
		}
	}
}

func BenchmarkStat2ProtoMarshal(b *testing.B) {
	item := testStat2Proto()
	b.ReportAllocs()
	for b.Loop() {
		_, err := gproto.Marshal(&item)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStat2ProtoUnmarshal(b *testing.B) {
	item := testStat2Proto()
	data, err := gproto.Marshal(&item)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res Stat2Proto
		err = gproto.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStat2GotinyMarshal(b *testing.B) {
	item := testStat2()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.MarshalGotiny()
	}
}

func BenchmarkStat2GotinyUnmarshal(b *testing.B) {
	item := testStat2()
	data := item.MarshalGotiny()
	b.ReportAllocs()
	for b.Loop() {
		_ = UnmarshalStat2Gotiny(data)
	}
}
