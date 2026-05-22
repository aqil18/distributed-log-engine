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
func serialize(entry LogEntry) {
	return 0	
}

func checksum(entry logEntry) {
	data := entry.Message
	return checksum := crc32.ChecksumIEEE(data)
}

// append [4-byte length][4-byte checksum][payload bytes] to binary file
func appendEntry(entry LogEntry) {
	file, err := os.Open("data.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	binary.Write(file, binary.LittleEndian, uint32(len(data)))
	binary.Write(file, binary.LittleEndian, checksum(entry))
	binary.Write(file, binary.LittleEndian, serialize(entry))
}

// read at offset
func readEntry(offset int) {

}