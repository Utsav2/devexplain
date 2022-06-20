package proc

import (
	ps "github.com/mitchellh/go-ps"
	"github.com/utsav2/devexp-monitor/lib"
)

type ProcInfo struct {
	Pid  int
	Name string
}

func ProcessInfos(cfg *lib.Config) (map[string][]ProcInfo, error) {
	ret := make(map[string][]ProcInfo)
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, proc := range processes {
		name := proc.Executable()
		targets, err := cfg.ProcessNameMatch(name)
		if err != nil {
			return nil, err
		}
		for _, t := range targets {
			target, ok := ret[t]
			if !ok {
				target = []ProcInfo{}
			}
			target = append(target, ProcInfo{
				Name: name,
				Pid:  proc.Pid(),
			})
			ret[t] = target
		}
	}
	return ret, nil
}
