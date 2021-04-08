package main

import (
	"bufio"
	"context"
	"github.com/djherbis/times"
	"github.com/fatih/color"
	"github.com/kasai-y/photo-dir-date/exiftool"
	"github.com/pkg/errors"
	"github.com/rodaine/table"
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

type resetFile struct {
	filepath string
	datetime time.Time
	createTs *time.Time
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
			Name:  "set-flat-time",
			Usage: "set the same datetime for the same reg-format photos. default, add the difference from the creation datetime of the first file.",
		},
		cli.BoolFlag{
			Name:     "dry-run",
			Usage:    "dry run.",
			Required: false,
			Hidden:   false,
		},
	}
	app.Action = action

	_ = app.Run(os.Args)
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
	if len(resetFiles) < 1 {
		println("no JPEG files exists.")
		return nil
	}

	//-- set datetime from filepath
	for fi, f := range resetFiles {
		sm := reg.FindStringSubmatch(f.filepath)
		if len(sm) < 7 {
			return nil
		}

		var y, m, d, h, i, s int

		for si, sv := range sm {
			if si == 0 {
				continue
			}
			v, err := strconv.Atoi(sv)
			if err != nil {
				return errors.WithStack(err)
			}

			switch si {
			case 1:
				y = v
			case 2:
				m = v
			case 3:
				d = v
			case 4:
				h = v
			case 5:
				i = v
			case 6:
				s = v
			}
		}
		resetFiles[fi].datetime = time.Date(y, time.Month(m), d, h, i, s, 0, time.Local)
	}

	//-- set createTs
	for i, f := range resetFiles {
		stat, err := times.Stat(f.filepath)
		if err != nil {
			return errors.WithStack(err)
		}
		if !stat.HasBirthTime() {
			continue
		}
		bt := stat.BirthTime()
		resetFiles[i].createTs = &bt
	}

	//-- add seconds from file created time.
	if !c.Bool("set-flat-time") {
		sameDateBaseTimeMap := make(map[string]time.Time)
		for i, f := range resetFiles {
			if f.createTs == nil {
				continue
			}
			k := f.datetime.String()
			if _, ok := sameDateBaseTimeMap[k]; !ok {
				sameDateBaseTimeMap[k] = *f.createTs
			}
			resetFiles[i].datetime = f.datetime.Add(f.createTs.Sub(sameDateBaseTimeMap[k]))
		}
	}

	//-- print table
	tbl := table.New("FilePath", "ModifiedTime", "CreateTime")
	tbl.WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithFirstColumnFormatter(color.New(color.FgCyan).SprintfFunc())
	for _, i := range resetFiles {
		tbl.AddRow(i.filepath, i.datetime.Format("2006-01-02 15:04:05.000000000 -07:00"), i.createTs.String())
	}
	tbl.Print()

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
			if err := exiftool.SetOriginalDateTime(i.filepath, i.datetime); err != nil {
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
