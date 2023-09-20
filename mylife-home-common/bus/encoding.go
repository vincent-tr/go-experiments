package bus

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

type encodingImpl struct{}

var Encoding = encodingImpl{}

func (e *encodingImpl) ReadString(buffer []byte) string {
	return string(buffer)
}

func (e *encodingImpl) WriteString(value string) []byte {
	return []byte(value)
}

func (e *encodingImpl) ReadBool(buffer []byte) bool {
	return buffer[0] != 0
}

func (e *encodingImpl) WriteBool(value bool) []byte {
	var iVal byte
	if value {
		iVal = 1
	} else {
		iVal = 0
	}

	return []byte{iVal}
}

func (e *encodingImpl) ReadUInt8(buffer []byte) uint8 {
	var data uint8
	e.read(buffer, data)
	return data
}

func (e *encodingImpl) WriteUInt8(value uint8) []byte {
	return e.write(value)
}

func (e *encodingImpl) ReadInt8(buffer []byte) int8 {
	var data int8
	e.read(buffer, data)
	return data
}

func (e *encodingImpl) WriteInt8(value int8) []byte {
	return e.write(value)
}

func (e *encodingImpl) ReadUInt32(buffer []byte) uint32 {
	var data uint32
	e.read(buffer, data)
	return data
}

func (e *encodingImpl) WriteUInt32(value uint32) []byte {
	return e.write(value)
}

func (e *encodingImpl) ReadInt32(buffer []byte) int32 {
	var data int32
	e.read(buffer, data)
	return data
}

func (e *encodingImpl) WriteInt32(value int32) []byte {
	return e.write(value)
}

func (e *encodingImpl) ReadFloat(buffer []byte) float64 {
	var data float64
	e.read(buffer, data)
	return data
}

func (e *encodingImpl) WriteFloat(value float64) []byte {
	return e.write(value)
}

func (e *encodingImpl) ReadJson(buffer []byte) any {
	var value any

	e.ReadTypedJson(buffer, &value)

	return value
}

func (e *encodingImpl) ReadTypedJson(buffer []byte, value any) {
	if err := json.Unmarshal(buffer, value); err != nil {
		panic(err)
	}
}

func (e *encodingImpl) WriteJson(value any) []byte {
	buf, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return buf
}

func (e *encodingImpl) read(buffer []byte, data any) {
	if err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &data); err != nil {
		panic(err)
	}
}

func (e *encodingImpl) write(data any) []byte {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		panic(err)
	}

	return buf.Bytes()
}
