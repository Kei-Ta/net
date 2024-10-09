package parse

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestParseEthernetFrame(t *testing.T) {
	// テスト用のデータを準備
	destMAC := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	srcMAC := net.HardwareAddr{0x11, 0x12, 0x13, 0x14, 0x15, 0x16}
	etherType := uint16(0x0800)               // IPv4を示すEtherType
	payload := []byte{0xde, 0xad, 0xbe, 0xef} // ペイロードのダミーデータ

	// バイトデータを構築
	data := append(destMAC, srcMAC...)
	etherTypeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(etherTypeBytes, etherType)
	data = append(data, etherTypeBytes...)
	data = append(data, payload...)

	// ParseEthernetFrame関数を呼び出す
	frame, err := ParseEthernetFrame(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 結果を検証
	if !reflect.DeepEqual(frame.Destination, destMAC) {
		t.Errorf("expected destination MAC %v, got %v", destMAC, frame.Destination)
	}

	if !reflect.DeepEqual(frame.Source, srcMAC) {
		t.Errorf("expected source MAC %v, got %v", srcMAC, frame.Source)
	}

	if frame.EtherType != etherType {
		t.Errorf("expected EtherType %v, got %v", etherType, frame.EtherType)
	}

	if !reflect.DeepEqual(frame.Payload, payload) {
		t.Errorf("expected payload %v, got %v", payload, frame.Payload)
	}
}

func TestParseEthernetFrameShortData(t *testing.T) {
	// データが短すぎる場合のテスト
	data := []byte{0x01, 0x02} // 長さが14バイト未満
	_, err := ParseEthernetFrame(data, false)
	if err == nil {
		t.Fatal("expected error for short data, but got none")
	}
	expectedError := "data too short for Ethernet frame"
	if err.Error() != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err.Error())
	}
}
