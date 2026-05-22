package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"log"
	"os"
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
func deserialize(buf []byte, messageLength uint32) LogEntry {
	var entry LogEntry
	var Timestamp int64
	var Level uint8
	var Message []byte
	err := binary.Read(buf, binary.LittleEndian, &Timestamp)
	err := binary.Read(buf, binary.LittleEndian, &Level)
	Message = buf[10:messageLength]

	var entry LogEntry
	LogEntry.Timestamp = Timestamp
	LogEntry.Level = Level
	LogEntry.Message = Message
	

	if err != nil {
		fmt.Println("Read failed:", err)
	}
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


