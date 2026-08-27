package server

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Session struct {
	db *sql.DB
	tx *sql.Tx
}

func (s *Session) handle(command, arg string) (string, error) {
	switch strings.ToUpper(command) {

	case "BEGIN":
		if s.tx != nil {
			return "", fmt.Errorf("already in a transaction")
		}
		tx, err := s.db.Begin()
		if err != nil {
			return "", err
		}
		s.tx = tx
		return "OK", nil

	case "COMMIT":
		if s.tx == nil {
			return "", fmt.Errorf("no transaction in progress")
		}
		err := s.tx.Commit()
		s.tx = nil
		if err != nil {
			return "", err
		}
		return "OK", nil

	case "ROLLBACK":
		if s.tx == nil {
			return "", fmt.Errorf("no transaction in progress")
		}
		err := s.tx.Rollback()
		s.tx = nil
		if err != nil {
			return "", err
		}
		return "OK", nil

	case "EXEC":
		if arg == "" {
			return "", fmt.Errorf("usage: EXEC <sql>")
		}
		var res sql.Result
		var err error
		if s.tx != nil {
			res, err = s.tx.Exec(arg)
		} else {
			res, err = s.db.Exec(arg)
		}
		if err != nil {
			return "", err
		}
		n, _ := res.RowsAffected()
		return fmt.Sprintf("OK %d rows affected", n), nil

	case "QUERY":
		if arg == "" {
			return "", fmt.Errorf("usage: QUERY <sql>")
		}
		var rows *sql.Rows
		var err error
		if s.tx != nil {
			rows, err = s.tx.Query(arg)
		} else {
			rows, err = s.db.Query(arg)
		}
		if err != nil {
			return "", err
		}
		defer rows.Close()
		return serializeRows(rows)

	default:
		return "", fmt.Errorf("unknown command: %s", command)
	}
}

func serializeRows(rows *sql.Rows) (string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(strings.Join(cols, "\t"))

	vals := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	rowCount := 0
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return "", err
		}
		strs := make([]string, len(cols))
		for i, v := range vals {
			strs[i] = fmt.Sprintf("%v", v)
		}
		sb.WriteString("\n" + strings.Join(strs, "\t"))
		rowCount++
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

type Server struct {
	db  *sql.DB
	sem chan struct{}
	wg  sync.WaitGroup
	ln  net.Listener
}

func NewServer(db *sql.DB, maxConns int) *Server {
	return &Server{
		db:  db,
		sem: make(chan struct{}, maxConns),
	}
}

func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.ln = ln
	log.Printf("listening on %s", addr)
	return s.serve()
}

func (s *Server) serve() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if isClosedErr(err) {
				return nil
			}
			log.Printf("accept error: %v", err)
			continue
		}

		select {
		case s.sem <- struct{}{}:
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				defer func() { <-s.sem }()
				s.handleConn(conn)
			}()
		default:
			conn.Write([]byte("ERR server busy\n"))
			conn.Close()
		}
	}
}

func isClosedErr(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr()
	log.Printf("connection opened: %s", remote)
	defer log.Printf("connection closed: %s", remote)

	sess := &Session{db: s.db}
	defer func() {
		if sess.tx != nil {
			_ = sess.tx.Rollback()
		}
	}()

	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		command := parts[0]
		arg := ""
		if len(parts) == 2 {
			arg = parts[1]
		}

		resp, err := sess.handle(command, arg)
		var out string
		if err != nil {
			out = "ERR " + err.Error()
		} else {
			out = resp
		}

		if _, err := conn.Write([]byte(out + "\n")); err != nil {
			return
		}
	}
}

func (s *Server) Shutdown() {
	if s.ln != nil {
		s.ln.Close()
	}
	s.wg.Wait()
}
