package zip

type Stats struct {
	Name, TmpFile    string
	UncompressedSize uint64
	CompressedSize   uint64
}

type StatsWriter interface {
	GetStats() Stats
}

func (self *fileWriter) GetStats() Stats {
	return Stats{
		Name:             self.header.Name,
		TmpFile:          self.tmp_filename,
		UncompressedSize: uint64(self.rawCount.count),
		CompressedSize:   uint64(self.compCount.count),
	}
}
