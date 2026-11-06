// Package server exposes the store over TCP using a tiny text protocol.
//
// Protocol (one command per line):
//
//	PING                -> PONG
//	SET <key> <value>   -> OK            (value may contain spaces)
//	GET <key>           -> <value> | (nil)
//	DEL <key>           -> 1 | 0
//	anything else       -> ERR <message>
package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/systemdr/distbuild/internal/store"
)

// Server accepts TCP connections and serves commands against a Store.
type Server struct {
	store *store.Store
	ln    net.Listener
	wg    sync.WaitGroup
}

// New creates a Server backed by s.
func New(s *store.Store) *Server { return &Server{store: s} }

// Listen binds to addr (e.g. ":7070").
func (srv *Server) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	return nil
}

// Addr returns the bound address (useful when listening on ":0" in tests).
func (srv *Server) Addr() string { return srv.ln.Addr().String() }

// Serve accepts connections until Close is called.
// Each client gets its own goroutine.
func (srv *Server) Serve() error {
	for {
		conn, err := srv.ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		srv.wg.Add(1)
		go func() {
			defer srv.wg.Done()
			srv.handle(conn)
		}()
	}
}

// Close stops accepting new connections and waits for active ones to finish.
func (srv *Server) Close() error {
	err := srv.ln.Close()
	srv.wg.Wait()
	return err
}

func (srv *Server) handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("read %s: %v", conn.RemoteAddr(), err)
			}
			return
		}
		reply := srv.execute(strings.TrimRight(line, "\r\n"))
		w.WriteString(reply + "\n")
		// Only flush when the client has no more pipelined commands waiting.
		if r.Buffered() == 0 {
			if err := w.Flush(); err != nil {
				return
			}
		}
	}
}

func (srv *Server) execute(line string) string {
	cmd, rest, _ := strings.Cut(line, " ")
	switch strings.ToUpper(cmd) {
	case "PING":
		return "PONG"
	case "SET":
		key, value, ok := strings.Cut(rest, " ")
		if !ok || key == "" {
			return "ERR usage: SET <key> <value>"
		}
		srv.store.Set(key, value)
		return "OK"
	case "GET":
		if rest == "" {
			return "ERR usage: GET <key>"
		}
		if v, ok := srv.store.Get(rest); ok {
			return v
		}
		return "(nil)"
	case "DEL":
		if rest == "" {
			return "ERR usage: DEL <key>"
		}
		if srv.store.Delete(rest) {
			return "1"
		}
		return "0"
	case "":
		return "ERR empty command"
	default:
		return fmt.Sprintf("ERR unknown command %q", cmd)
	}
}
