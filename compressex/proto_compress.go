package compressex

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

func ProtoMarshal(v any) ([]byte, error) {
	dataBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buff, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(dataBytes)
	if err != nil {
		return nil, err
	}
	_ = writer.Close()

	return buff.Bytes(), nil
}

func ProtoUnmarshal(dataBytes []byte, v any) error {
	if len(dataBytes) == 0 {
		return nil
	}
	reader, err := gzip.NewReader(bytes.NewReader(dataBytes))
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	_, err = io.Copy(&buff, reader)
	if err != nil {
		return err
	}

	return json.Unmarshal(buff.Bytes(), v)
}
