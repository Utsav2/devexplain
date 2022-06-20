package proc

import "github.com/shirou/gopsutil/v3/process"

func RssMem(pid int) (uint64, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, err
	}
	mem, err := proc.MemoryInfo()
	if err != nil {
		return 0, err
	}
	return mem.RSS, nil
}
