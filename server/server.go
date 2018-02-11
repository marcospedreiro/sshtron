package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/marcospedreiro/sshtron/config"
	"github.com/marcospedreiro/sshtron/game"
	"github.com/marcospedreiro/sshtron/player"
	"golang.org/x/crypto/ssh"
)

func handleConnection(newConn net.Conn, gm *game.GameManager, sshConfig *ssh.ServerConfig) {
	// first perform handshake on incoming net.Conn
	sshConn, chans, reqs, err := ssh.NewServerConn(newConn, sshConfig)
	if err != nil {
		fmt.Println("Failed to handshake with new client")
	}

	// service incoming request channel
	go ssh.DiscardRequests(reqs)

	// service incoming channel
	for newChannel := range chans {
		/* channels have a type depending on application level
		protocol intended. In the case of a shell, the type is
		"session" and ServerShell may be used to present a simple
		terminal interface */
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "channel type is not session, rejecting")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Println("could not accept channel.")
			return
		}
		fmt.Println("new player joined")
		// To see how many concurrent users are online
		//fmt.Printf("Player joined. Current stats: %d users, %d games\n", gm.SessionCount(), gm.GameCount())

		// Reject all out of band requests and accept the unix defaults, pty-req and shell
		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "pty-req":
					req.Reply(true, nil)
					continue
				case "shell":
					req.Reply(true, nil)
					continue
				}
				req.Reply(false, nil)
			}
		}(requests)
		gm.HandleNewChannel(channel, sshConn.User())
	}

	return
}

func setupSSHServer(cfg *config.Config) (net.Listener, *ssh.ServerConfig) {
	// let anyone login
	sshConf := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateKeyPath := fmt.Sprintf("%s%s", *cfg.Server.SSHKeyPath, *cfg.Server.SSHKeyName)
	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to load private key")
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Failed to parse private key")
	}
	sshConf.AddHostKey(privateKey)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", *cfg.Server.SSHPort))
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("failed to listen for connection")
	}

	return listener, sshConf
}

// Start setups a game server
func Start(cfg *config.Config) error {

	// setup listeners
	go func() {
		httpPort := fmt.Sprintf("0.0.0.0:%s", *cfg.Server.HTTPPort)
		panic(http.ListenAndServe(httpPort, http.FileServer(http.Dir(*cfg.Server.HTTPFileServerDir))))
	}()

	sshListener, sshConfig := setupSSHServer(cfg)

	fmt.Printf("HTTP listener on port: %s\nSSH listener on port: %s\n",
		*cfg.Server.HTTPPort,
		*cfg.Server.SSHPort,
	)

	// Set game properties defined in the configuration file
	SetConfigurableProperties(cfg)

	// Create the GameManager
	gm := game.NewGameManager()

	for {
		newConn, err := sshListener.Accept()
		if err != nil {
			panic("failed to accept incoming connection")
		}

		go handleConnection(newConn, gm, sshConfig)
	}
}

// SetConfigurableProperties invokes package specific property setters that overrides
// default or sets optional values based on what is in the configuration object
func SetConfigurableProperties(cfg *config.Config) {
	game.SetGameManagerProperties(cfg)
	game.SetGameServerProperties(cfg)
	player.SetPlayerProperties(cfg)
	return
}
