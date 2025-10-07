package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/leinonen/bbs/bbs"
	"github.com/leinonen/bbs/config"
	"github.com/leinonen/bbs/ui"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHServer struct {
	config   *config.Config
	db       *sql.DB
	listener net.Listener
	sessions *bbs.SessionManager
}

func NewSSHServer(cfg *config.Config, db *sql.DB) *SSHServer {
	return &SSHServer{
		config:   cfg,
		db:       db,
		sessions: bbs.NewSessionManager(),
	}
}

func (s *SSHServer) Start() error {
	sshConfig := &ssh.ServerConfig{
		NoClientAuth: s.config.AllowAnonymous,
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			user, err := bbs.AuthenticateUser(s.db, conn.User(), string(password))
			if err != nil {
				return nil, fmt.Errorf("invalid credentials")
			}
			return &ssh.Permissions{
				Extensions: map[string]string{
					"user-id": fmt.Sprintf("%d", user.ID),
				},
			}, nil
		},
	}

	hostKey, err := s.loadOrGenerateHostKey()
	if err != nil {
		return fmt.Errorf("failed to load host key: %v", err)
	}
	sshConfig.AddHostKey(hostKey)

	listener, err := net.Listen("tcp", s.config.ListenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			if s.listener == nil {
				return nil
			}
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn, sshConfig)
	}
}

func (s *SSHServer) Stop() {
	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

func (s *SSHServer) handleConnection(netConn net.Conn, config *ssh.ServerConfig) {
	defer netConn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(netConn, config)
	if err != nil {
		log.Printf("Failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()

	log.Printf("New SSH connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("Could not accept channel: %v", err)
			continue
		}

		go s.handleSession(channel, requests, sshConn)
	}
}

func (s *SSHServer) handleSession(channel ssh.Channel, requests <-chan *ssh.Request, sshConn *ssh.ServerConn) {
	defer channel.Close()

	term := terminal.NewTerminal(channel, "")

	var user *bbs.User
	if s.config.AllowAnonymous && sshConn.Permissions == nil {
		user = &bbs.User{
			Username: "anonymous",
			ID:       0,
		}
	} else if sshConn.Permissions != nil {
		userID := sshConn.Permissions.Extensions["user-id"]
		user, _ = bbs.GetUserByID(s.db, userID)
	}

	session := s.sessions.CreateSession(user, term)
	defer s.sessions.RemoveSession(session.ID)

	go func() {
		for req := range requests {
			switch req.Type {
			case "pty-req":
				termLen := req.Payload[3]
				term.SetPrompt("")
				width, height := parseDims(req.Payload[termLen+4:])
				term.SetSize(width, height)
				req.Reply(true, nil)
			case "shell":
				req.Reply(true, nil)
			case "window-change":
				width, height := parseDims(req.Payload)
				term.SetSize(width, height)
			default:
				req.Reply(false, nil)
			}
		}
	}()

	ui := ui.NewUI(term, s.db, session)
	ui.Run()
}

func (s *SSHServer) loadOrGenerateHostKey() (ssh.Signer, error) {
	keyPath := s.config.HostKeyPath

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return s.generateHostKey(keyPath)
	}

	privateBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(privateBytes)
}

func (s *SSHServer) generateHostKey(path string) (ssh.Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return nil, err
	}

	return ssh.NewSignerFromKey(key)
}

func parseDims(b []byte) (int, int) {
	if len(b) < 8 {
		return 80, 24
	}
	width := int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
	height := int(b[4])<<24 | int(b[5])<<16 | int(b[6])<<8 | int(b[7])
	return width, height
}
