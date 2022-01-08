package utils

import (
	"bufio"
	"math/rand"
	"net"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/roj1512/chife/config"
)

var possibilities = strings.Split("0123456789abcdef", "")

func EnsurePastesDirectoryExistence() {
	_, err := os.Stat(config.PASTES_DIRECTORY)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		err = os.Mkdir(config.PASTES_DIRECTORY, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func NewSlug() string {
	rand.Seed(time.Now().UnixMicro())
	slug := ""

	for i := 0; i < 6; i++ {
		slug += possibilities[rand.Intn(len(possibilities))]
	}

	return slug
}

func GetFileName(slug string) string {
	return path.Join(config.PASTES_DIRECTORY, slug+".txt")
}

func Read(reader *bufio.Reader) []rune {
	runes := []rune{}

	for {
		if len(runes) > config.MAX_LEN {
			break
		}

		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}

		runes = append(runes, r)
	}

	return runes
}

func Respond(conn net.Conn, writer *bufio.Writer, result string) error {
	_, err := writer.WriteString(result + "\n")
	if err != nil {
		log.Error("Writing to connection: " + err.Error())
		return err
	}

	err = writer.Flush()
	if err != nil {
		log.Error("Flushing writer: " + err.Error())
		return err
	}

	err = conn.Close()
	if err != nil {
		log.Error("Closing connection: " + err.Error())
		return err
	}

	return nil
}

func RespondWithError(conn net.Conn, writer *bufio.Writer, err error) {
	Respond(conn, writer, "error: "+err.Error())
}
