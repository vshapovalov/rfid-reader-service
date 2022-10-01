package readers

func calcCRCChecksum(data []byte) byte {
	reminder := 0
	for i := 0; i < len(data); i++ {
		reminder = (reminder + int(data[i])) % 256
	}
	return byte(reminder)
}
