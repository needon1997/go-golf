package golf_test

import (
	"fmt"
	"github.com/needon1997/go-golf"
	"testing"
)

func TestName(t *testing.T) {
	golf.App.RegisterDaemonGo(func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
	})
	golf.App.Run()
}
