package test

import (
	"net/http"
	"testing"
)

func BenchmarkRank(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:8090/rank")
		b.Log(resp, err)
	}
}
