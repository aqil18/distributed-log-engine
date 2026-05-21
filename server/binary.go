type LogEntry struct {
    Timestamp int64
    Level     uint8
    Message   []byte
}

//Level represents the severity of the log message — e.g. debug, info, warn, error.
// It's typically mapped to constants:
const (
    DEBUG uint8 = iota  // 0
    INFO                 // 1
    WARN                 // 2
    ERROR                // 3
)

// serialize bytes of log entry
func serialize() {
	data := entry.Message
	checksum := crc32.ChecksumIEEE(data)
}

// append [4-byte length][4-byte checksum][payload bytes] to binary file
func appendEntry(entry LogEntry) {
	binary.Write(w, binary.LittleEndian, uint32(len(data)))
	binary.Write(w, binary.LittleEndian, checksum)
	w.Write(data)
}

// read at offset
func readEntry(offset int) {

}