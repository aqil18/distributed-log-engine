package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp int64
	Level     uint8
	Message   []byte
}

// Level represents the severity of the log message — e.g. debug, info, warn, error.
const (
	DEBUG uint8 = iota // 0
	INFO               // 1
	WARN               // 2
	ERROR              // 3
)

// we must serialize as file.Write() requires a []byte
func serialize(entry LogEntry) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, entry.Timestamp)
	binary.Write(buf, binary.LittleEndian, entry.Level)
	buf.Write(entry.Message)
	return buf.Bytes()
}

// deserialize into logentry
func deserialize(data []byte) LogEntry {
    r := bytes.NewReader(data)
    var entry LogEntry
    binary.Read(r, binary.LittleEndian, &entry.Timestamp) // reads 8 bytes
    binary.Read(r, binary.LittleEndian, &entry.Level)     // reads 1 byte
    entry.Message = data[9:]                              // everything after
    return entry
}

func checksum(entry LogEntry) uint32 {
	return crc32.ChecksumIEEE(serialize(entry))
}

// append [4-byte length][4-byte checksum][payload bytes] to binary file
func appendEntry(entry LogEntry) {
	file, err := os.OpenFile("logs.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := serialize(entry)
	binary.Write(file, binary.LittleEndian, uint32(len(data)))
	binary.Write(file, binary.LittleEndian, checksum(entry))
	file.Write(data)
}


func writeEntry(message string, level uint8) error {
	entry := LogEntry{
		Timestamp: time.Now().UnixNano(),
		Level:     level,
		Message:   []byte(message),
	}
	appendEntry(entry)
	return nil
}

// read at offset
func readEntry(offset int64) (LogEntry, error) {
    file, err := os.OpenFile("logs.bin", os.O_RDONLY, 0644)
    if err != nil {
        return LogEntry{}, err
    }
    defer file.Close()

    file.Seek(offset, io.SeekStart) // jump to byte position

    var length uint32
    binary.Read(file, binary.LittleEndian, &length)  // read 4 bytes

    var storedChecksum uint32
    binary.Read(file, binary.LittleEndian, &storedChecksum)  // read 4 bytes

    data := make([]byte, length)
    io.ReadFull(file, data)  // read exactly length bytes

    if crc32.ChecksumIEEE(data) != storedChecksum {
        return LogEntry{}, fmt.Errorf("checksum mismatch, entry corrupted")
    }

    return deserialize(data), nil
}



func main() {
	writeEntry("This is an important log", 2)
	writeEntry("This is not so important and is more informational", 1)
	data, err := readEntry(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data.Message))
	data, err = readEntry(41)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data.Message))

}