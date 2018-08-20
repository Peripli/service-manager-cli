package util

import (
	"bufio"
	"errors"
	"io"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadInput(input io.ReadWriter) (string, error) {
	bufReader := bufio.NewReader(input)
	line, isPrefix, err := bufReader.ReadLine()
	if isPrefix {
		return "", errors.New("input too long")
	}
	if err != nil {
		return "", err
	}

	return string(line), nil
}

func ReadPassword() (string, error) {
	password, err := terminal.ReadPassword((int)(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(password), nil
}
