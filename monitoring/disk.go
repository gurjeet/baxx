package monitoring

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"syscall"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/jinzhu/gorm"
	"github.com/mackerelio/go-osstat/disk"
	"github.com/mackerelio/go-osstat/memory"
	log "github.com/sirupsen/logrus"
)

type DiskUsagePerNode struct {
	ID     uint64 `gorm:"primary_key"`
	Kind   string `gorm:"type:varchar(255) not null"`
	NodeID string `gorm:"type:varchar(255) not null"`

	DiskAll  uint64
	DiskUsed uint64
	DiskFree uint64

	CreatedAt time.Time
}

type DiskIOPerNode struct {
	ID     uint64 `gorm:"primary_key"`
	Kind   string `gorm:"type:varchar(255) not null"`
	NodeID string `gorm:"type:varchar(255) not null"`

	DiskReadsCompleted  uint64
	DiskWritesCompleted uint64

	CreatedAt time.Time
}

type DiskMDPerNode struct {
	ID     uint64 `gorm:"primary_key"`
	Kind   string `gorm:"type:varchar(255) not null"`
	NodeID string `gorm:"type:varchar(255) not null"`

	ExitCode int
	MDADM    string `gorm:"type:text not null"`

	CreatedAt time.Time
}

type MemStatsPerNode struct {
	ID     uint64 `gorm:"primary_key"`
	NodeID string `gorm:"type:varchar(255) not null"`

	MemoryTotal      uint64
	MemoryUsed       uint64
	MemoryBuffers    uint64
	MemoryCached     uint64
	MemoryFree       uint64
	MemoryAvailable  uint64
	MemoryActive     uint64
	MemoryInactive   uint64
	MemorySwapTotal  uint64
	MemorySwapUsed   uint64
	MemorySwapCached uint64
	MemorySwapFree   uint64
	CreatedAt        time.Time
}

func GetMDADM(md string) DiskMDPerNode {
	out := DiskMDPerNode{}
	out.NodeID = Hostname()
	out.Kind = md

	cmd := exec.Command("mdadm", "-D", "-t", "/dev/"+md)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			log.Panic(err)
		}
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	out.ExitCode = exitCode
	out.MDADM = fmt.Sprintf("%s%s", outStr, errStr)
	return out
}

func GetDiskUsage(path string) DiskUsagePerNode {
	disk := DiskUsagePerNode{Kind: path, NodeID: Hostname()}
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		log.Panic(err)
	}
	disk.DiskAll = fs.Blocks * uint64(fs.Bsize)
	disk.DiskFree = fs.Bfree * uint64(fs.Bsize)
	disk.DiskUsed = disk.DiskAll - disk.DiskFree

	return disk
}

func GetDiskIOStats(diskName string) DiskIOPerNode {
	out := DiskIOPerNode{}
	out.NodeID = Hostname()
	out.Kind = diskName

	disksStats, err := disk.Get()
	var diskIO disk.Stats
	if err != nil {
		log.Panic(err)
	} else {
		found := false
		for _, d := range disksStats {
			if d.Name == diskName {
				diskIO = d
				found = true
				break
			}
		}
		if !found {
			log.Panicf("missing disk %s", diskName)
		}
	}

	out.DiskReadsCompleted = diskIO.ReadsCompleted
	out.DiskWritesCompleted = diskIO.WritesCompleted
	return out
}

func GetMemoryStats() MemStatsPerNode {
	m, err := memory.Get()
	if err != nil {
		log.Panic(err)
	}

	return MemStatsPerNode{
		NodeID:           Hostname(),
		MemoryTotal:      m.Total,
		MemoryUsed:       m.Used,
		MemoryBuffers:    m.Buffers,
		MemoryCached:     m.Cached,
		MemoryFree:       m.Free,
		MemoryAvailable:  m.Available,
		MemoryActive:     m.Active,
		MemoryInactive:   m.Inactive,
		MemorySwapTotal:  m.SwapTotal,
		MemorySwapUsed:   m.SwapUsed,
		MemorySwapCached: m.SwapCached,
		MemorySwapFree:   m.SwapFree,
	}
}

type Stats struct {
	Mem    []MemStatsPerNode
	DU     []DiskUsagePerNode
	IO     []DiskIOPerNode
	MDADM  []DiskMDPerNode
	Watch  []MonitoringPerNode
	NodeID string
}

func GetStats(db *gorm.DB, node_id string, n int) ([]*Stats, error) {
	watch := []MonitoringPerNode{}
	if err := db.Order("id").Find(&watch).Error; err != nil {
		return nil, err
	}
	nm := map[string]bool{}

	for _, w := range watch {
		nm[w.NodeID] = true
	}
	out := []*Stats{}
	nodes := []string{}
	for n := range nm {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	for _, nodeId := range nodes {
		mem := []MemStatsPerNode{}
		du := []DiskUsagePerNode{}
		dio := []DiskIOPerNode{}
		mdadm := []DiskMDPerNode{}

		if err := db.Limit(n).Where("node_id = ?", nodeId).Order("id desc").Find(&mem).Error; err != nil {
			return nil, err
		}
		if err := db.Limit(n).Where("node_id = ?", nodeId).Order("id desc").Find(&du).Error; err != nil {
			return nil, err
		}
		if err := db.Limit(n).Where("node_id = ?", nodeId).Order("id desc").Find(&dio).Error; err != nil {
			return nil, err
		}
		if err := db.Limit(n).Where("node_id = ?", nodeId).Order("id desc").Find(&mdadm).Error; err != nil {
			return nil, err
		}

		st := &Stats{
			Mem:    mem,
			DU:     du,
			IO:     dio,
			MDADM:  mdadm,
			NodeID: nodeId,
			Watch:  []MonitoringPerNode{},
		}
		for _, w := range watch {
			if w.NodeID == nodeId {
				st.Watch = append(st.Watch, w)
			}
		}
		out = append(out, st)
	}
	return out, nil
}

func reverse(numbers []float64) []float64 {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
func (s *Stats) ASCII(w io.Writer, height int) {
	data := []float64{}

	if len(s.Watch) > 0 {
		fmt.Fprintf(w, "%s Who Watches The Watchers:\n", s.NodeID)
		for _, m := range s.Watch {
			fmt.Fprintf(w, "%s %s tick: %s (%.2fs) ago, schedule: %.0fs\n", m.NodeID, m.Kind, m.Tick.Format(time.ANSIC), time.Since(m.Tick).Seconds(), m.Schedule)
		}
		fmt.Fprintf(w, "\n")
	}

	for _, m := range s.Mem {
		data = append(data, float64(m.MemoryFree)/float64(1024*1024*1024))
	}
	if len(data) > 0 {
		data = reverse(data)
		fmt.Fprintf(w, "%s\n\n", asciigraph.Plot(data, asciigraph.Height(height), asciigraph.Caption(fmt.Sprintf("%s Free Memory: %fGB", s.NodeID, data[len(data)-1]))))
	}

	data = []float64{}
	for _, m := range s.DU {
		data = append(data, float64(m.DiskUsed)/float64(1024*1024*1024))
	}
	if len(data) > 0 {
		data = reverse(data)
		fmt.Fprintf(w, "%s\n\n", asciigraph.Plot(data, asciigraph.Height(height), asciigraph.Caption(fmt.Sprintf("%s Used Disk: %fGB", s.NodeID, data[len(data)-1]))))
	}

	if len(s.MDADM) > 0 {
		fmt.Fprintf(w, "%s:\n%s\n\n", s.NodeID, s.MDADM[0].MDADM)
	}
}
