package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// BridgeStream is a binary serialization type
type BridgeStream struct {
	buffer     *bytes.Buffer
	readIndex  int
	writeIndex int
}

// BridgeSerializer is a template for custom serializer types
type BridgeSerializer interface {
	Read(stream *BridgeStream)
	Write(stream *BridgeStream)
}

// WriteInt adds the given integer to the next slot in the stream
func (stream *BridgeStream) WriteInt(data int) error {
	fixedInt := int32(data)
	err := binary.Write(stream.buffer, binary.LittleEndian, &fixedInt)
	if err != nil {
		fmt.Println("Could not write the integer:", err)
	}
	stream.readIndex += 4
	return err
}

// ReadInt returns the next integer value from buffer
func (stream *BridgeStream) ReadInt() (n int, err error) {
	var data int32
	err = binary.Read(stream.buffer, binary.LittleEndian, &data)
	if err != nil {
		fmt.Println("Could not read the integer:", err)
		return
	}

	stream.readIndex += 4
	return int(data), err
}

// WriteIntArray adds the given integer array to the next slot in the stream
func (stream *BridgeStream) WriteIntArray(data []int) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the integer array:", err)
		return err
	}
	for i, datum := range data {
		err := stream.WriteInt(datum)
		if err != nil {
			fmt.Printf("Could not read the value %d index of the integer array: %s", i, err)
			return err
		}
	}
	return err
}

// ReadIntArray returns the next integer array from the stream
func (stream *BridgeStream) ReadIntArray() (data []int, err error) {
	length, err := stream.ReadInt()
	if err != nil {
		fmt.Println("Could not read the length of the integer array:", err)
		return nil, err
	}
	raw := make([]int, length)
	for i := range raw {
		datum, err := stream.ReadInt()
		if err != nil {
			fmt.Printf("Could not read the value %d index of the integer array: %s", i, err)
			return nil, err
		}
		raw[i] = datum
	}
	return raw, nil
}

// WriteFloat adds the given float to the next slot in the stream
func (stream *BridgeStream) WriteFloat(data float32) error {
	err := binary.Write(stream.buffer, binary.LittleEndian, &data)
	if err != nil {
		fmt.Println("Could not write the float:", err)
	}

	stream.readIndex += 4
	return err
}

// ReadFloat returns the next float value from buffer
func (stream *BridgeStream) ReadFloat() (n float32, err error) {
	var data float32
	err = binary.Read(stream.buffer, binary.LittleEndian, &data)
	if err != nil {
		fmt.Println("Could not read the float:", err)
	}

	stream.readIndex += 4
	return data, err
}

// WriteFloatArray adds the given float array to the next slot in the stream
func (stream *BridgeStream) WriteFloatArray(data []float32) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the float array:", err)
		return err
	}
	for i, datum := range data {
		err := stream.WriteFloat(datum)
		if err != nil {
			fmt.Printf("Could not read the value %d index of the float array: %s", i, err)
			return err
		}
	}
	return err
}

// ReadFloatArray returns the next float array from the stream
func (stream *BridgeStream) ReadFloatArray() (data []float32, err error) {
	length, err := stream.ReadInt()
	if err != nil {
		fmt.Println("Could not read the length of the float array:", err)
		return nil, err
	}
	raw := make([]float32, length)
	for i := range raw {
		datum, err := stream.ReadFloat()
		if err != nil {
			fmt.Printf("Could not read the value %d index of the float array: %s", i, err)
			return nil, err
		}
		raw[i] = datum
	}
	return raw, nil
}

// WriteBool adds the given boolean to the next slot in the stream
func (stream *BridgeStream) WriteBool(data bool) error {
	err := binary.Write(stream.buffer, binary.LittleEndian, &data)
	if err != nil {
		fmt.Println("Could not write the integer:", err)
	}

	stream.readIndex += 1
	return err
}

// ReadBool returns the next boolean value from buffer
func (stream *BridgeStream) ReadBool() (n bool, err error) {
	var data bool
	err = binary.Read(stream.buffer, binary.LittleEndian, &data)

	stream.readIndex += 1
	return data, err
}

// WriteBoolArray adds the given bool array to the next slot in the stream
func (stream *BridgeStream) WriteBoolArray(data []bool) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the bool array:", err)
		return err
	}
	for i, datum := range data {
		err := stream.WriteBool(datum)
		if err != nil {
			fmt.Printf("Could not read the value %d index of the bool array: %s", i, err)
			return err
		}
	}
	return err
}

// ReadBoolArray returns the next bool array from the stream
func (stream *BridgeStream) ReadBoolArray() (data []bool, err error) {
	length, err := stream.ReadInt()
	if err != nil {
		fmt.Println("Could not read the length of the bool array:", err)
		return nil, err
	}
	raw := make([]bool, length)
	for i := range raw {
		datum, err := stream.ReadBool()
		if err != nil {
			fmt.Printf("Could not read the value %d index of the bool array: %s", i, err)
			return nil, err
		}
		raw[i] = datum
	}
	return raw, nil
}

// WriteString adds the given string to the next slot in the stream
func (stream *BridgeStream) WriteString(data string) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the string:", err)
		return err
	}
	nBytes, err := stream.buffer.WriteString(data)
	if err != nil {
		fmt.Println("Could not write the bytearray of string:", err)
	}
	stream.readIndex += nBytes
	return err
}

// ReadString returns the next string value from buffer
func (stream *BridgeStream) ReadString() (data string, err error) {
	length, err := stream.ReadInt()
	if err != nil {
		fmt.Println("Could not read the length of the string:", err)
		return "", err
	}
	raw := make([]byte, length)
	err = binary.Read(stream.buffer, binary.LittleEndian, &raw)
	stream.readIndex += length
	return string(raw), err
}

// WriteStringArray adds the given string array to the next slot in the stream
func (stream *BridgeStream) WriteStringArray(data []string) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the string array:", err)
		return err
	}
	for i, datum := range data {
		err := stream.WriteString(datum)
		if err != nil {
			fmt.Printf("Could not read the value %d index of the string array: %s", i, err)
			return err
		}
	}
	return err
}

// ReadStringArray returns the next string array from the stream
func (stream *BridgeStream) ReadStringArray() (data []string, err error) {
	length, err := stream.ReadInt()
	if err != nil {
		fmt.Println("Could not read the length of the string array:", err)
		return nil, err
	}
	raw := make([]string, length)
	for i := range raw {
		datum, err := stream.ReadString()
		if err != nil {
			fmt.Printf("Could not read the value %d index of the string array: %s", i, err)
			return nil, err
		}
		raw[i] = datum
	}

	return raw, nil
}

// WriteBytes adds the given byte array to the next slot in the stream
func (stream *BridgeStream) WriteBytes(data []byte) error {
	err := stream.WriteInt(len(data))
	if err != nil {
		fmt.Println("Could not write the length of the string:", err)
		return err
	}

	nBytes, err := stream.buffer.Write(data)
	if err != nil {
		fmt.Println("Could not write the bytearray of string:", err)
	}
	stream.writeIndex += nBytes
	return err
}

// ReadBytes returns the next byte array from stream
func (stream *BridgeStream) ReadBytes() (data []byte, err error) {
	length, _ := stream.ReadInt()
	data = make([]byte, length)
	err = binary.Read(stream.buffer, binary.LittleEndian, &data)

	stream.writeIndex += length
	return data, err
}

// WriteStream adds the given stream to the next slot in the stream
func (stream *BridgeStream) WriteStream(data *BridgeStream) error {
	err := stream.WriteBytes(data.Encode())
	if err != nil {
		fmt.Println("Could not write the substream:", err)
	}
	return err
}

// ReadStream returns the next sub stream within the stream
func (stream *BridgeStream) ReadStream() (data *BridgeStream, err error) {
	raw, err := stream.ReadBytes()
	if err != nil {
		fmt.Println("Could not read the substream:", err)
		return nil, err
	}
	return &BridgeStream{
		buffer:     bytes.NewBuffer(raw),
		readIndex:  0,
		writeIndex: 0,
	}, err
}

// WriteStream adds the given stream to the next slot in the stream
func (stream *BridgeStream) Write(serializer BridgeSerializer) error {
	serializerStream := New()
	serializer.Write(&serializerStream)
	stream.WriteStream(&serializerStream)
	return nil
}

// ReadStream returns the next sub stream within the stream
func (stream *BridgeStream) Read(serializer BridgeSerializer) error {
	serializerStream, err := stream.ReadStream()
	if err != nil {
		fmt.Println("Could not read the serializer substream:", err)
		return err
	}
	serializer.Read(serializerStream)
	return nil
}

// Clear clean ups all data and resets the state of the stream
func (stream *BridgeStream) Clear() {
	stream.buffer.Reset()
	stream.readIndex = 0
	stream.writeIndex = 0
}

// Encode returns the stream as a bytearray
func (stream *BridgeStream) Encode() []byte {
	return stream.buffer.Bytes()
}

// Empty reports the emptiness of the stream
func (stream *BridgeStream) Empty() bool {
	return stream.buffer.Len() == 0
}

// HasMore returns true if the stream has more data to read
func (stream *BridgeStream) HasMore() bool {
	return stream.buffer.Len() != 0
}

// New creates an empty BridgeStream
func New() BridgeStream {
	stream := BridgeStream{
		buffer:     new(bytes.Buffer),
		readIndex:  0,
		writeIndex: 0,
	}
	stream.Clear()
	return stream
}

func main(){}
