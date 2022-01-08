package chife

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/roj1512/chife/config"
	"github.com/roj1512/chife/utils"
	log "github.com/sirupsen/logrus"
)

func Start() {
	initialize()

	listener, err := net.Listen("tcp", config.ADDRESS)
	if err != nil {
		panic(err)
	}

	log.Info("Chife is running at " + config.ADDRESS + "...")

	for {
		conn, err := listener.Accept()
		go handleConn(conn, err)
	}
}

func initialize() {
	utils.EnsurePastesDirectoryExistence()
}

func handleConn(conn net.Conn, err error) {
	if err != nil {
		log.Error("Accepting connection: " + err.Error())
		return
	}

	conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	data := string(utils.Read(reader))

	if len(strings.TrimSpace(data)) == 0 {
		log.Warn("Empty input")
		utils.RespondWithError(conn, writer, errors.New("empty input"))
		return
	}

	slug := utils.NewSlug()
	err = os.WriteFile(utils.GetFileName(slug), []byte(data), os.ModePerm)

	if err != nil {
		log.Warn("Writing to file: " + err.Error())
		utils.RespondWithError(conn, writer, err)
		return
	}

	utils.Respond(conn, writer, slug)
}
