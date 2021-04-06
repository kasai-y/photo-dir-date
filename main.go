package main

import (
	"bufio"
	"context"
	"github.com/kasai-y/photo-dir-date/exiftool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type dateStruct struct {
	y int
	m int
	d int
	h int
	i int
	s int
	n int
}

func (d *dateStruct) time() time.Time {
	return time.Date(d.y, time.Month(d.m), d.d, d.h, d.i, d.s, d.n, time.Local)
}

type resetFile struct {
	filepath string
	datetime dateStruct
}

func main() {

	log.Logger = log.With().Caller().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	app := cli.NewApp()
	app.Name = "photo-dir-d"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "d,dir",
			Usage:    "target directory",
			Required: true,
		},
		cli.StringFlag{
			Name:     "reg-format",
			Usage:    "datetime format as regexp",
			Required: false,
			Hidden:   true,
			Value:    "(\\d{4})-(\\d{2})-(\\d{2})[^\\d]*(\\d{2})(\\d{2})(\\d{2})",
		},
		cli.BoolFlag{
			Name:     "dry-run",
			Usage:    "dry run.",
			Required: false,
			Hidden:   false,
		},
	}
	app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Error().Stack().Err(err).Send()
	}
}

func action(c *cli.Context) {

	if err := exiftool.Init(); err != nil {
		if err == exiftool.ErrNoExiftool {
			println("exiftool is not found.")
			os.Exit(0)
		}
		log.Error().Stack().Err(err).Send()
		os.Exit(-1)
	}

	if err := run(c); err != nil {
		log.Error().Stack().Err(err).Send()
		os.Exit(-1)
	}

	os.Exit(0)
}

func run(c *cli.Context) error {

	dir := c.String("dir")
	format := c.String("reg-format")
	reg, err := regexp.Compile(format)
	if err != nil {
		return errors.WithStack(err)
	}

	var resetFiles []resetFile

	//-- fide files
	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Name()[0] == '.' {
			return nil
		}
		if !reg.MatchString(path) {
			return nil
		}
		if ext := strings.ToLower(filepath.Ext(path)); ext != ".jpg" && ext != ".jpeg" {
			return nil
		}

		resetFiles = append(resetFiles, resetFile{
			filepath: path,
		})
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	//-- set datetime
	sameDateMap := make(map[string]int)
	for i, f := range resetFiles {
		sm := reg.FindStringSubmatch(f.filepath)
		if len(sm) < 7 {
			return nil
		}

		for si, sv := range sm {
			if si == 0 {
				continue
			}
			if strings.Trim(sv, "_") == "" {
				continue
			}
			v, err := strconv.Atoi(sv)
			if err != nil {
				return errors.WithStack(err)
			}

			switch si {
			case 1:
				resetFiles[i].datetime.y = v
			case 2:
				resetFiles[i].datetime.m = v
			case 3:
				resetFiles[i].datetime.d = v
			case 4:
				resetFiles[i].datetime.h = v
			case 5:
				resetFiles[i].datetime.i = v
			case 6:
				resetFiles[i].datetime.s = v
			}
		}
		dt := resetFiles[i].datetime

		layout := "2006-01-02 15:04:05"
		if _, ok := sameDateMap[dt.time().Format(layout)]; !ok {
			sameDateMap[dt.time().Format(layout)] = 0
		} else {
			sameDateMap[dt.time().Format(layout)]++
		}
		resetFiles[i].datetime.n = sameDateMap[dt.time().Format(layout)]
	}

	println()
	for _, i := range resetFiles {
		println(i.datetime.time().Format(exiftool.DateTimeOriginalLayout) + " | " + i.filepath)
	}

	if c.Bool("dry-run") {
		println("-- dry-run --")
		return nil
	}

	if !confirm() {
		return nil
	}

	println("continue to set DateTime.")

	//--- run
	eg, _ := errgroup.WithContext(context.Background())
	ch := make(chan bool, 10)
	for _, i := range resetFiles {
		ch <- true
		i := i
		eg.Go(func() error {
			defer func() { <-ch }()
			if err := exiftool.SetOriginalDateTime(i.filepath, i.datetime.time()); err != nil {
				return errors.WithStack(err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func confirm() bool {
	println("continue? y/n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "y":
			return true
		case "n":
			return false
		}
	}
	return false
}
