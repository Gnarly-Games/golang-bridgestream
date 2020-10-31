package main

import (
	"reflect"
	"testing"
)

type MatchInfo struct {
	matchID     int
	playerIDs   []int
	playerNames []string
}

func (matchInfo *MatchInfo) Write(stream *BridgeStream) {
	stream.WriteInt(matchInfo.matchID)
	stream.WriteIntArray(matchInfo.playerIDs)
	stream.WriteStringArray(matchInfo.playerNames)
}
func (matchInfo *MatchInfo) Read(stream *BridgeStream) {
	id, _ := stream.ReadInt()
	matchInfo.matchID = id
	matchInfo.playerIDs, _ = stream.ReadIntArray()
	matchInfo.playerNames, _ = stream.ReadStringArray()
}

func TestBridgeStream(t *testing.T) {

	t.Run("Integer I/O", func(t *testing.T) {
		var expected int = 1231
		stream := New()
		stream.WriteInt(expected)
		got, _ := stream.ReadInt()

		if got != expected {
			t.Errorf("got %d want %d", got, expected)
		}
	})

	t.Run("Integer Collection I/O", func(t *testing.T) {
		expected := []int{123, 456}
		stream := New()
		stream.WriteIntArray(expected)
		got, _ := stream.ReadIntArray()

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v want %v", got, expected)
		}
	})

	t.Run("String I/O", func(t *testing.T) {
		var expected string = "1231"
		stream := New()
		stream.WriteString(expected)
		got, _ := stream.ReadString()

		if got != expected {
			t.Errorf("got %s want %s", got, expected)
		}
	})

	t.Run("String Collection I/O", func(t *testing.T) {
		expected := []string{"456", "123"}
		stream := New()
		stream.WriteStringArray(expected)
		got, _ := stream.ReadStringArray()

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v want %v", got, expected)
		}
	})

	t.Run("Bytes I/O", func(t *testing.T) {
		expected := []byte("1231")
		stream := New()
		stream.WriteBytes(expected)
		got, _ := stream.ReadBytes()

		if string(got) != string(expected) {
			t.Errorf("got %s want %s", got, expected)
		}
	})

	t.Run("Float I/O", func(t *testing.T) {
		var expected float32 = 0.123
		stream := New()
		stream.WriteFloat(expected)
		got, _ := stream.ReadFloat()

		if got != expected {
			t.Errorf("got %f want %f", got, expected)
		}
	})

	t.Run("Float Collection I/O", func(t *testing.T) {
		expected := []float32{0.2, 0.5}
		stream := New()
		stream.WriteFloatArray(expected)
		got, _ := stream.ReadFloatArray()

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v want %v", got, expected)
		}
	})

	t.Run("Bool I/O", func(t *testing.T) {
		expected := true
		stream := New()
		stream.WriteBool(expected)
		got, _ := stream.ReadBool()

		if got != expected {
			t.Errorf("got %t want %t", got, expected)
		}
	})

	t.Run("Bool Collection I/O", func(t *testing.T) {
		expected := []bool{true, false, true}
		stream := New()
		stream.WriteBoolArray(expected)
		got, _ := stream.ReadBoolArray()

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("got %v want %v", got, expected)
		}
	})

	t.Run("Stream I/O", func(t *testing.T) {
		expected := true
		substream := New()
		substream.WriteBool(expected)
		stream := New()
		stream.WriteStream(&substream)

		receivedSubstream, _ := stream.ReadStream()
		got, _ := receivedSubstream.ReadBool()

		if got != expected {
			t.Errorf("got %v want %v", got, expected)
		}
	})

	t.Run("Serializer I/O", func(t *testing.T) {

		expectedMatch := &MatchInfo{
			matchID:     1231,
			playerIDs:   []int{1, 2},
			playerNames: []string{"rufus", "dufus"},
		}
		stream := New()

		stream.Write(expectedMatch)

		receivedMatch := &MatchInfo{}
		stream.Read(receivedMatch)

		if !reflect.DeepEqual(expectedMatch, receivedMatch) {
			t.Errorf("got %v want %v", receivedMatch, expectedMatch)
		}
	})
	t.Run("Clear", func(t *testing.T) {

		stream := New()

		stream.WriteInt(123)

		stream.Clear()

		if !reflect.DeepEqual(stream.Encode(), []byte{}) {
			t.Errorf("Stream is not empty: %v", stream.Encode())
		}
	})

	t.Run("Empty / Not Empty", func(t *testing.T) {

		stream := New()

		if !stream.Empty() {
			t.Errorf("Stream is not empty: %v", stream.Encode())
		}

		stream.WriteInt(123)

		if stream.Empty() {
			t.Errorf("Stream should not be empty: %v", stream.Encode())
		}
	})

	t.Run("Has More", func(t *testing.T) {

		stream := New()

		stream.WriteInt(123)
		stream.WriteString("test")

		if !stream.HasMore() {
			t.Errorf("Stream is expected to have more data: %v", stream.Encode())
		}

		stream.ReadInt()
		if !stream.HasMore() {
			t.Errorf("Stream is expected to have more data: %v", stream.Encode())
		}

		stream.ReadString()
		if stream.HasMore() {
			t.Errorf("Stream should not have any data: %v", stream.Encode())
		}
	})
}
