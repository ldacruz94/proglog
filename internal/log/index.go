package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offwidth = 4
	posWidth = 8
	entWidth = offwidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())

	if err := os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}

	// memory-map the index file so we can read and write to it like a byte slice.
	// It is faster than traditional file I/O and allows us to
	// treat file contents as if they were in memory.
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}

	return idx, nil
}

func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	if in == -1 {
		out = uint32((i.size / uint64(entWidth)) - 1)
	} else {
		out = uint32(in)
	}

	pos = uint64(out) * uint64(entWidth)
	if i.size < pos+uint64(entWidth) {
		return 0, 0, io.EOF
	}

	out = enc.Uint32(i.mmap[pos : pos+uint64(offwidth)])
	pos = enc.Uint64(i.mmap[pos+uint64(offwidth) : pos+uint64(entWidth)])
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+uint64(entWidth) {
		return io.EOF
	}

	enc.PutUint32(i.mmap[i.size:i.size+uint64(offwidth)], off)
	enc.PutUint64(i.mmap[i.size+uint64(offwidth):i.size+uint64(entWidth)], pos)
	i.size += uint64(entWidth)
	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Close(); err != nil {
		return err
	}
	return nil
}
