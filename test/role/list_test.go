package role

//testing
import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	//testing
	//go test -bench=.
	//go test --timeout 9999999999999s
)

func TestMain(u *testing.T) {
	__(u)

	select {}
}

func __(u *testing.T) {
	fmt.Printf("\033[1;32m%s\033[0m\n", strings.ReplaceAll(u.Name(), "Test", ""))
}

func Benchmark1(u *testing.B) {
	u.ReportAllocs()

	for u.Loop() {

	}
}

func Benchmark2(u *testing.B) {
	u.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

		}
	})
}

func cmd(name string, v ...string) {
	c := exec.Command(name, v...)
	r, err := c.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(r))
}
