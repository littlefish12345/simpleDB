package simpleDB

import (
	"os"
	"bytes"
	"errors"
	"encoding/binary"
)

var (
	file *os.File = nil
	DBDamaged error = errors.New("Database has been damaged.")
	DBKeyNotFound error = errors.New("Database key not found.")
	DBNotOpened error = errors.New("Database not opened yet.")
)

func CreateDatabase(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(bytes.Repeat([]byte{0}, 18))
	return nil
}

func OpenDatabase(path string) error {
	var err error
	file, err = os.OpenFile(path, os.O_RDWR, 6)
	return err
}

func CloseDatabase() {
	file.Close()
	file = nil
}

func WriteDatabase(key int64, value int64) error {
	if file == nil {
		return DBNotOpened
	}
	var pointer int64 = 0
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	fileLength := fi.Size()
	for i := 0; i < 63; i++ {
		tempPointer, err := trackPointer(key, i, pointer)
		if err != nil {
			if err != DBKeyNotFound {
				return err
			} else {
				file.WriteAt(bytes.Repeat([]byte{0}, 18), fileLength)
				bit := getBit(key, i)
				if bit {
					file.WriteAt([]byte{0x01}, pointer + 9)
					file.WriteAt(convertInt64ToBytes(fileLength), pointer + 10)
				} else {
					file.WriteAt([]byte{0x01}, pointer)
					file.WriteAt(convertInt64ToBytes(fileLength), pointer + 1)
				}
				pointer = fileLength
				fileLength = fileLength + 18
			}
		} else {
			pointer = tempPointer
		}
	}
	bit := getBit(key, 63)
	if bit {
		file.WriteAt([]byte{0x01}, pointer + 9)
		file.WriteAt(convertInt64ToBytes(value), pointer + 10)
	} else {
		file.WriteAt([]byte{0x01}, pointer)
		file.WriteAt(convertInt64ToBytes(value), pointer + 1)
	}

	return nil
}

func ReadDatabase(key int64) (int64, error) {
	if file == nil {
		return 0, DBNotOpened
	}

	var pointer int64 = 0
	var err error
	for i := 0; i < 64; i++ {
		pointer, err = trackPointer(key, i, pointer)
		if err != nil {
			return 0, err
		}
	}
	return pointer, err
}

func trackPointer(key int64, bitNum int, pointer int64) (int64, error) {
	buffer1 := make([]byte, 1)
	buffer2 := make([]byte, 8)
	bit := getBit(key, bitNum)
	if bit {
		pointer = pointer + 9
	}

	length, err := file.ReadAt(buffer1, pointer)
	if err != nil {
		return 0, err
	}
	if length != 1 {
		return 0, DBDamaged
	}

	if buffer1[0] == 0x00 {
		return 0, DBKeyNotFound
	}

	length, err = file.ReadAt(buffer2, pointer + 1)
	if err != nil {
		return 0, err
	}
	if length != 8 {
		return 0, DBDamaged
	}

	return convertBytesToInt64(buffer2), nil
}

func convertBytesToInt64(data []byte) int64 {
	var temp int64
	bytesBuffer := bytes.NewBuffer(data[:8])
	binary.Read(bytesBuffer, binary.BigEndian, &temp)
	return int64(temp)
}

func convertInt64ToBytes(data int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &data)
	return bytesBuffer.Bytes()[:8]
}

func getBit(num int64, bit int) bool {
	if (num >> bit) & 0x01 == 1 {
		return true
	}
	return false
}
