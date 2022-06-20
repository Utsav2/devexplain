package lib

import (
	"fmt"

	"github.com/utsav2/devexplain/lib/sink"
	"github.com/utsav2/devexplain/lib/sink/multi"
	"go.starlark.net/starlark"
)

type Config struct {
	thread  *starlark.Thread
	sources map[string]starlark.Value
	Sink    sink.Sink
}

func New(progname string, configs []string) (*Config, error) {
	cfg := &Config{
		thread:  &starlark.Thread{},
		sources: map[string]starlark.Value{},
	}

	allSinks := []string{}
	for _, c := range configs {
		globals, err := starlark.ExecFile(cfg.thread, c, nil, nil)
		if err != nil {
			return nil, err
		}
		sinkFunc, ok := globals["sinks"]
		if ok {
			if sinkFunc.Type() != "function" {
				return nil, fmt.Errorf("invalid config %s: sinks is not a function", c)
			}
			sinks, err := starlark.Call(cfg.thread, sinkFunc, nil, nil)
			if err != nil {
				return nil, fmt.Errorf("invalid config %s: sink function errored out: %s", c, err)
			}
			if sinks.Type() != "dict" {
				return nil, fmt.Errorf("invalid config %s: sink return is not a dictionary", c)
			}
			sinksDict := sinks.(*starlark.Dict)
			for _, k := range sinksDict.Keys() {
				if k.Type() != "string" {
					return nil, fmt.Errorf("invalid config %s: sinks dictionary keys are not strings", c)
				}
				topCfg, found, err := sinksDict.Get(k)
				if err != nil {
					return nil, fmt.Errorf("invalid config %s: sinks dictionary key: err retrieving value", err)
				}
				if !found {
					return nil, fmt.Errorf("invalid config %s: sinks dictionary key: dictionary values changed while iterating", k)
				}
				if topCfg.Type() != "dict" {
					return nil, fmt.Errorf("invalid config %s: sinks dictionary value: value not dictionary", k)
				}
				// eventually we might pass these dictionaries down to the sinks for initialization.
				allSinks = append(allSinks, k.(starlark.String).GoString())
			}
		}

		sourcesFunc, ok := globals["sources"]
		if !ok {
			return nil, fmt.Errorf("invalid config %s: no global source specified", c)
		}
		if sourcesFunc.Type() != "function" {
			return nil, fmt.Errorf("invalid config %s: sources is not a function", c)
		}
		sources, err := starlark.Call(cfg.thread, sourcesFunc, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("invalid config %s: sources function errored out: %s", c, err)
		}
		if sources.Type() != "dict" {
			return nil, fmt.Errorf("invalid config %s: sources is not a dictionary", c)
		}
		d := sources.(*starlark.Dict)
		for _, k := range d.Keys() {
			if k.Type() != "string" {
				return nil, fmt.Errorf("invalid config %s: sources dictionary keys are not strings", c)
			}
			topCfg, found, err := d.Get(k)
			if err != nil {
				return nil, fmt.Errorf("invalid config %s: sources dictionary key: err retrieving value", err)
			}
			if !found {
				return nil, fmt.Errorf("invalid config %s: sources dictionary key: dictionary values changed while iterating", k)
			}
			if topCfg.Type() != "dict" {
				return nil, fmt.Errorf("invalid config %s: sources dictionary value: value not dictionary", k)
			}
			process, found, err := topCfg.(*starlark.Dict).Get(starlark.String("process"))
			if err != nil {
				return nil, fmt.Errorf("invalid config %s: sources dictionary value: err retrieving value: %s", c, err)
			}
			if !found {
				return nil, fmt.Errorf("invalid config %s: sources dictionary value: key 'process' not found", c)
			}
			if process.Type() != "dict" {
				return nil, fmt.Errorf("invalid config %s: process config: value not dictionary", c)
			}
			val, found, err := process.(*starlark.Dict).Get(starlark.String("name"))
			if err != nil {
				return nil, fmt.Errorf("invalid config %s: process dictionary value: err retrieving value: %s", c, err)
			}
			if !found {
				return nil, fmt.Errorf("invalid config %s: process dictionary value: key 'name' not found", c)
			}
			switch val.Type() {
			case "string":
			case "function":
			default:
				return nil, fmt.Errorf("invalid config %s: sources value type has to be string, list or func, got %s", c, val.Type())
			}
			key := k.(starlark.String).GoString()
			// we overwrite target configs from newer sources.
			cfg.sources[key] = val
		}
	}
	if len(allSinks) == 0 {
		allSinks = append(allSinks, "stderr")
	}
	sink, err := multi.New(progname, allSinks)
	if err != nil {
		return nil, err
	}
	cfg.Sink = sink
	return cfg, nil
}

func match(thread *starlark.Thread, val starlark.Value, name string) (bool, error) {
	switch val.Type() {
	case "string":
		return val.(starlark.String).GoString() == name, nil
	case "function":
		res, err := starlark.Call(thread, val, starlark.Tuple{starlark.String(name)}, nil)
		if err != nil {
			return false, fmt.Errorf("error matching %s with config: %s", name, err)
		}
		if res.Type() != "bool" {
			return false, fmt.Errorf("error matching %s starlark function returned non bool value: %s", name, res)
		}
		ret, _ := res.(starlark.Bool)
		return ret == starlark.True, nil
	default:
		return false, fmt.Errorf("internal error: type %s not handled in config parsing", val.Type())
	}
}

func (c *Config) ProcessNameMatch(n string) ([]string, error) {
	matches := []string{}
	for t, val := range c.sources {
		matched, err := match(c.thread, val, n)
		if err != nil {
			return nil, err
		}
		if matched {
			matches = append(matches, t)
		}
	}
	return matches, nil
}

func (c *Config) Sources() []string {
	sources := []string{}
	for k := range c.sources {
		sources = append(sources, k)
	}
	return sources
}
