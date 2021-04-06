package exiftool

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
	"time"
)

var exiftool string

var ErrNoExiftool = errors.New("no exiftool exists.")

func Init() error {
	output, err := exec.Command("which", "exiftool").Output()
	if err != nil {
		return errors.WithStack(err)
	}
	exiftool = strings.Trim(string(output), "\n")

	if exiftool == "" {
		return ErrNoExiftool
	}

	return nil
}

const DateTimeOriginalLayout = "2006-01-02 15:04:05.000 -07:00"

func SetOriginalDateTime(file string, ts time.Time) error {

	cmd := exec.Command(
		exiftool,
		"-overwrite_original",
		"-DateTimeOriginal=\""+ts.Format(DateTimeOriginalLayout)+"\"",
		fmt.Sprintf("-SubSecTimeOriginal=%d", ts.Nanosecond()),
		file,
	)

	output, err := cmd.Output()
	if err != nil {
		return errors.WithStack(err)
	}
	println(cmd.String() + "\n" + string(output))

	return nil
}
