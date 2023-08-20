package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"tailscale.com/tsnet"
)

func main() {
	s := &tsnet.Server{}
	if dir, ok := os.LookupEnv("TS_STATE_DIR"); ok {
		s.Dir = dir
	}
	if hostname, ok := os.LookupEnv("TS_HOSTNAME"); ok {
		s.Hostname = hostname
	} else if hostname, err := os.Hostname(); err == nil {
		s.Hostname = hostname
	}
	defer s.Close()

	ln, err := s.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	lc, err := s.LocalClient()
	if err != nil {
		log.Fatal(err)
	}

	serveChan := make(chan net.Conn, 10)
	for i := 0; i < 50; i++ {
		go func() {
			for r := range serveChan {
				who, err := lc.WhoIs(context.TODO(), r.RemoteAddr().String())
				if err != nil {
					r.Close()
				}

				if repo, err := os.Open("/museum/repo.fossil"); err == nil {
					repo.Close()
				} else if os.IsNotExist(err) {
					log.Printf("serve: creating new repo for %s", who.UserProfile.LoginName)
					if err := createRepo(who.UserProfile.LoginName); err != nil {
						log.Printf("serve: %v", err)
					}
				} else {
					log.Print(err)
				}

				cmd := exec.Cmd{
					Path:   "/bin/fossil",
					Args:   []string{"/bin/fossil", "http", "/museum/repo.fossil", "--jsmode", "bundled"},
					Env:    []string{"PATH=/bin", fmt.Sprintf("REMOTE_USER=%s", who.UserProfile.LoginName)},
					Stdin:  r,
					Stdout: r,
					Stderr: os.Stderr,
				}
				if err := cmd.Run(); err != nil {
					log.Printf("serve: %v", err)
				}
			}
		}()
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		serveChan <- conn
	}
}

func createRepo(adminUser string) error {
	cmd := exec.Cmd{
		Path:   "/bin/fossil",
		Args:   []string{"/bin/fossil", "new", "--admin-user", adminUser, "/museum/repo.fossil"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create repo: %w", err)
	}

	cmd = exec.Cmd{
		Path:   "/bin/fossil",
		Args:   []string{"/bin/fossil", "sql", "-R", "/museum/repo.fossil", "INSERT OR REPLACE INTO config VALUES ('remote_user_ok',1,strftime('%s','now'));"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create repo: enable remote_user_ok: %w", err)
	}

	return nil
}
