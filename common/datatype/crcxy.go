package datatype

func CheckSum(at_str string) []byte {
	var y uint8
	for i := 0; i < len(at_str); i += 2 {
		y += BytesToUInt8(HexStringToByte(at_str[i : i+2]))
	}
	return UInt8ToBytes(y)
}
