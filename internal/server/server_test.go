package server

import (
	"bufio"
	"net"
	"testing"

	"github.com/systemdr/distbuild/internal/store"
)

func TestProtocol(t *testing.T) {
	srv := New(store.New())
	if err := srv.Listen("127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}
	go srv.Serve()
	defer srv.Close()

	conn, err := net.Dial("tcp", srv.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r := bufio.NewReader(conn)

	cases := []struct{ send, want string }{
		{"PING", "PONG"},
		{"GET user:1", "(nil)"},
		{"SET user:1 Ada Lovelace", "OK"},
		{"GET user:1", "Ada Lovelace"},
		{"set user:1 Grace", "OK"},
		{"GET user:1", "Grace"},
		{"DEL user:1", "1"},
		{"DEL user:1", "0"},
		{"SET onlykey", "ERR usage: SET <key> <value>"},
		{"FLY away", `ERR unknown command "FLY"`},
	}
	for _, c := range cases {
		if _, err := conn.Write([]byte(c.send + "\n")); err != nil {
			t.Fatal(err)
		}
		got, err := r.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		if got = got[:len(got)-1]; got != c.want {
			t.Errorf("%q: got %q want %q", c.send, got, c.want)
		}
	}
}
