package aggregator

import (
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type FileHandlers struct {
	heapfile *os.File
	hmfile   *os.File
}

type FileHandlersInterface interface {
	GetHeapFile() (*os.File, error)
	GetHmFile() (*os.File, error)
	WriteHk(hk Hashkey) error
	ReadHm() (Hashmap, error)
	Write(data []byte) (int64, error)
	WriteDeletion() (int64, error)
	Retrieve(start int64) ([]byte, error)
	CreateDbIfNotExists() error
}

type Hashkey struct {
	Key   string
	Start int64
}

type Hashmap []Hashkey

var mountpath = os.Getenv("MOUNTPATH")
var hmFileName = os.Getenv("HM_BACKUP")
var heapFileName = os.Getenv("HEAP_FILE")
var hmBackup = fmt.Sprintf("%s/%s", mountpath, hmFileName)
var heapfile = fmt.Sprintf("%s/%s", mountpath, heapFileName)

func (hm Hashmap) GetStart(key string) (int64, bool) {
	for i := len(hm) - 1; i >= 0; i-- {
		if hm[i].Key == key {
			return hm[i].Start, true
		}
	}
	return 0, false
}

func (fh *FileHandlers) GetHeapFile() (*os.File, error) {
	if fh.heapfile == nil {
		f, err := os.OpenFile(heapfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		fh.heapfile = f
	}
	return fh.heapfile, nil
}

func (fh *FileHandlers) GetHmFile() (*os.File, error) {
	if fh.hmfile == nil {
		f, err := os.OpenFile(hmBackup, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
		fh.hmfile = f
	}
	return fh.hmfile, nil
}

func (fh *FileHandlers) WriteHk(hk Hashkey) error {
	f, err := fh.GetHmFile()
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("%d,%s\n", hk.Start, hk.Key))
	if err != nil {
		return err
	}
	return f.Sync()
}

func (fh *FileHandlers) ReadHm() (Hashmap, error) {
	hm := Hashmap{}
	f, err := fh.GetHmFile()
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(f)
	hks, err := csvReader.ReadAll()
	if err != nil {
		return hm, err
	}
	for _, hk := range hks {
		if len(hk) != 2 {
			return hm, fmt.Errorf("malformed Hashkey received %v", hk)
		}
		start, err := strconv.Atoi(hk[0])
		if err != nil {
			return hm, err
		}
		hm = append(hm, Hashkey{
			Start: int64(start),
			Key:   hk[1],
		})
	}
	return hm, nil
}

func (fh *FileHandlers) Write(data []byte) (int64, error) {
	if len(data) > 32767 {
		return 0, errors.New("maximum data size exceeded (32767 bytes)")
	}
	f, err := fh.GetHeapFile()
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	start := fi.Size()

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(len(data)))
	b = append(b, data...)

	_, err = f.Write(b)
	if err != nil {
		return 0, err
	}

	return start, f.Sync()
}

func (fh *FileHandlers) WriteDeletion() (int64, error) {
	f, err := fh.GetHeapFile()
	if err != nil {
		return 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	start := fi.Size()
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(uint16(1)<<15))
	_, err = f.Write(b)
	if err != nil {
		return 0, err
	}

	return start, f.Sync()
}

func (fh *FileHandlers) Retrieve(start int64) ([]byte, error) {
	f, err := fh.GetHeapFile()
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(start, 0); err != nil {
		return nil, err
	}
	b := make([]byte, 2)
	if _, err := f.Read(b); err != nil {
		return nil, err
	}
	val := binary.LittleEndian.Uint16(b)
	// file has tombstone - has been deleted
	if val>>15 == 1 {
		return nil, nil
	}
	size := val &^ 0b1000000000000000
	b = make([]byte, size)
	if _, err := f.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (fh *FileHandlers) CreateDbIfNotExists() error {
	_, err := os.Stat(mountpath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(mountpath, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = os.Stat(heapfile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(heapfile); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = os.Stat(hmBackup)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if _, err := os.Create(hmBackup); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
