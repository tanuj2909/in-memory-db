package file

import (
	"fmt"
	"time"

	"github.com/tanuj2909/in-memory-db/app/store"
)

func checkHeader(data []byte) ([]byte, error) {
	if len(data) < 9 {
		return data, fmt.Errorf("not enough data: expected atleast 9 bytes, got %d", len(data))
	}

	if string(data[:5]) != "REDIS" {
		return data, fmt.Errorf("invalid header provided: %s", string(data[:5]))
	}

	version := string(data[5:9])
	fmt.Printf("REDIS file version: %s\n", version)
	return data[9:], nil
}

func readLength(data []byte) (int, bool, []byte, error) {

	if len(data) == 0 {
		return 0, false, data, fmt.Errorf("no data to read")
	}

	firstByte := data[0]

	msb := firstByte >> 6

	switch msb {
	case 0:
		// length is present in the first byte inself
		return readIgnoringTwoMSB(firstByte), false, data[1:], nil
	case 1:
		// length is present iun current 6 bits + next byte
		if len(data) < 2 {
			return 0, false, data, fmt.Errorf("not enough  data to read")
		}
		return readIntIgnoringTwoMSB(data[:2]), false, data[2:], nil
	case 2:
		// length is present in the next 4 bytes
		if len(data) < 5 {
			return 0, false, data, fmt.Errorf("not enough data to read")
		}
		return int(readAsInteger(data[1:5])), false, data[5:], nil
	case 3:
		lsb := firstByte & 0b00111111
		var bytesToRead []byte
		var nextBytes []byte

		switch lsb {
		case 0:
			// 8 bit integer
			if len(data) < 2 {
				return 0, false, data, fmt.Errorf("not enough data to read")
			}
			bytesToRead = data[1:3]
			nextBytes = data[2:]
		case 1:
			// 16 bit integer
			if len(data) < 3 {
				return 0, false, data, fmt.Errorf("not enough data to read")
			}
			bytesToRead = data[1:3]
			nextBytes = data[3:]
		case 2:
			// 32 bit integer
			if len(data) < 5 {
				return 0, false, data, fmt.Errorf("not enough data to read")
			}
			bytesToRead = data[1:5]
			nextBytes = data[5:]
		default:
			return 0, false, data, fmt.Errorf("unimplemented value for LSB")
		}

		return int(readAsInteger(bytesToRead)), true, nextBytes, nil
	default:
		return 0, false, data, fmt.Errorf("invalid valuse for MSB")
	}

}

func readString(data []byte) (string, []byte, error) {
	n, special, reamainingData, err := readLength(data)
	if err != nil {
		return "", data, fmt.Errorf("error readign string length: %v", err)
	}

	if special {
		return fmt.Sprintf("%d", n), reamainingData, nil
	}

	if len(reamainingData) < n {
		return "", data, fmt.Errorf("not enough data to read")
	}

	return string(reamainingData[:n]), reamainingData[n:], nil
}
func readAuxiliaryField(data []byte) (string, string, []byte, error) {
	if len(data) == 0 {
		return "", "", data, fmt.Errorf("unexpected end of auxiliary section")
	}

	key, remainingData, err := readString(data)
	if err != nil {
		return "", "", data, fmt.Errorf("error reading key: %v", err)
	}

	value, remainingData, err := readString(remainingData)
	if err != nil {
		return "", "", data, fmt.Errorf("error reading value: %v", err)
	}

	return key, value, remainingData, nil
}
func parseFile(data []byte) error {

	//header section
	data, err := checkHeader(data)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("no data to read after header in the file")
	}

	//auxiliary section
	for data[0] == 0xFA {
		key, value, remainingData, err := readAuxiliaryField(data[1:])
		if err != nil {
			return fmt.Errorf("error reading auxliary field: %v", err)
		}
		fmt.Printf("auxiliary field: %s -> %s\n", key, value)
		data = remainingData
	}

	//database section
	if len(data) == 0 {
		return fmt.Errorf(("file data ended abruptly after reading auxiliary fields"))
	}
	if data[0] == 0xFE {
		if len(data) < 2 {
			return fmt.Errorf("not enough data to read while reading database section")
		}
		fmt.Printf("database selector: %d\n", data[1])
		data = data[2:]
	}

	// resize DB fields
	if len(data) == 0 {
		return fmt.Errorf("file data ended abruptly after reading database selector")
	}

	if data[0] == 0xFB {
		hashSize, remainingData, err := readInteger(data[1:])
		if err != nil {
			return fmt.Errorf("error reading hash size: %v", err)
		}
		fmt.Printf("Hash Size: %d\n", hashSize)
		expireHashSize, remainingData, err := readInteger(remainingData)
		if err != nil {
			return fmt.Errorf("error reading expire hash size: %v", err)
		}
		fmt.Printf("Expire Hash Size: %d\n", expireHashSize)
		data = remainingData
	}

	// read the key value pairs
	for len(data) > 0 && data[0] != 0xFF {
		remainingData, err := processKeyValuePair(data)
		if err != nil {
			return fmt.Errorf("error reading key value pair: %v", err)
		}
		data = remainingData
	}

	return nil
}

func readInteger(data []byte) (int, []byte, error) {
	n, _, remainingData, err := readLength(data)
	if err != nil {
		return 0, data, fmt.Errorf("error reading integer: %v", err)
	}
	return n, remainingData, nil
}

func readExpiry(data []byte) (int64, []byte, error) {
	if len(data) == 0 {
		return 0, data, fmt.Errorf("unexpected end of data when parsing for (optional) expiry")
	}

	switch data[0] {
	case 0xFD:
		// Expiry timestamp is present in seconds
		if len(data) < 5 {
			return 0, data, fmt.Errorf("not enough data to read: expected atleast 5 bytes, got %d", len(data))
		}
		exp := readAsInteger(data[1:5])
		return exp * 1000, data[5:], nil

	case 0xFC:
		// Expiry timestamp is present in milliseconds
		if len(data) < 9 {
			return 0, data, fmt.Errorf("not enough data to read: expected atleast 9 bytes, got %d", len(data))
		}
		return readAsInteger(data[1:9]), data[9:], nil

	default:
		// No expiry timestamp present
		return -1, data, nil
	}
}

func processKeyValuePair(data []byte) ([]byte, error) {

	if len(data) == 0 {
		return data, fmt.Errorf("unexpected end of data while reading key value pairs")
	}

	expiry, remainingData, err := readExpiry(data)
	if err != nil {
		return data, fmt.Errorf("error reading expiry: %w", err)
	}

	key, remainingData, err := readString(remainingData[1:])
	if err != nil {
		return data, fmt.Errorf("error reading key: %w", err)
	}

	value, remainingData, err := readString(remainingData)
	if err != nil {
		return data, fmt.Errorf("error reading value: %w", err)
	}

	if expiry == -1 || expiry > time.Now().UnixMilli() {
		store.Store.Set(key, value, expiry)
	}

	return remainingData, nil
}
