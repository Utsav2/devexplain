package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/utsav2/devexplain/lib"
	"github.com/utsav2/devexplain/lib/proc"
)

var progname = "devexplain"

func _monitor(cfg *lib.Config) error {
	procInfos, err := proc.ProcessInfos(cfg)
	if err != nil {
		return err
	}
	for target, infos := range procInfos {
		for i, p := range infos {
			mem, err := proc.RssMem(p.Pid)
			if err != nil {
				return err
			}
			cfg.Sink.Metric(fmt.Sprintf("%s_%d", target, i), float64(mem))
		}
	}
	return nil
}

func monitor() error {
	cfgs, err := findConfigs(flag.Args())
	if err != nil {
		return err
	}
	if len(cfgs) == 0 {
		return fmt.Errorf(
			"no config file found. Either specify one in args, or in a .%src.star in the enclosing directory or a parent directory. See https://github.com/utsav2/%s/blob/main/examples/ for examples",
			progname,
			progname,
		)
	}
	cfg, err := lib.New(progname, cfgs)
	if err != nil {
		return err
	}
	cfg.Sink.Log("initialized: %s", cfg.Sources())
	err = _monitor(cfg)
	if err != nil {
		cfg.Sink.Err(err)
	}
	cfg.Sink.Log("run complete")
	return nil
}

func findConfigs(args []string) ([]string, error) {
	cfgs := make([]string, 0)
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		potentialCfg := strings.Join([]string{p, ".", progname, "rc", ".star"}, "")
		stat, err := os.Stat(potentialCfg)
		if err == nil && !stat.IsDir() {
			cfgs = append(cfgs, potentialCfg)
		}
		dir := filepath.Dir(p)
		if dir == p {
			break
		}
		p = dir
	}
	cfgs = append(cfgs, args...)
	return cfgs, nil
}

func run(cadence int64) error {
	dur := time.Duration(cadence) * time.Minute
	log.Printf("cadence: %v\n", dur)

	// run it once so users get quick feedback.
	if err := monitor(); err != nil {
		return err
	}
	if dur > 0 {
		for range time.Tick(dur) {
			if err := monitor(); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}

func main() {
	// this is a startup option as we don't want a checked in config file to be able to change this.
	cadence := flag.Int64(
		"cadence_minutes",
		-1,
		"how often (in minutes) to scrape and export data. Pass a negative number to run once and quit.",
	)
	flag.Parse()
	if err := run(*cadence); err != nil {
		log.Fatal(err)
	}
}
