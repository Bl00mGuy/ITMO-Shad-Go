//go:build !solution

package jsonlist

import (
	"bufio"
	"encoding/json"
	"io"
	"reflect"
)

func Marshal(writer io.Writer, slice interface{}) error {
	if err := validateSlice(slice); err != nil {
		return err
	}

	serializer := newSliceSerializer(writer, slice)

	return serializer.serialize()
}

func Unmarshal(reader io.Reader, slice interface{}) error {
	if err := validateSlicePointer(slice); err != nil {
		return err
	}

	deserializer := newSliceDeserializer(reader, slice)

	return deserializer.deserialize()
}

func validateSlice(slice interface{}) error {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	return nil
}

func validateSlicePointer(slice interface{}) error {
	if reflect.TypeOf(slice).Kind() != reflect.Ptr {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	if reflect.ValueOf(slice).Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: reflect.TypeOf(slice)}
	}
	return nil
}

type SliceSerializer struct {
	writer   io.Writer
	slice    interface{}
	encoder  *json.Encoder
	buffered *bufio.Writer
}

func newSliceSerializer(writer io.Writer, slice interface{}) *SliceSerializer {
	bufferedWriter := bufio.NewWriter(writer)
	encoder := json.NewEncoder(bufferedWriter)

	return &SliceSerializer{
		writer:   writer,
		slice:    slice,
		encoder:  encoder,
		buffered: bufferedWriter,
	}
}

func (s *SliceSerializer) serialize() error {
	sliceValue := reflect.ValueOf(s.slice)

	for i := 0; i < sliceValue.Len(); i++ {
		if err := s.encodeElement(sliceValue.Index(i).Interface()); err != nil {
			return err
		}
	}

	if err := s.flush(); err != nil {
		return err
	}

	return nil
}

func (s *SliceSerializer) encodeElement(element interface{}) error {
	return s.encoder.Encode(element)
}

func (s *SliceSerializer) flush() error {
	return s.buffered.Flush()
}

type SliceDeserializer struct {
	reader  io.Reader
	slice   interface{}
	decoder *json.Decoder
}

func newSliceDeserializer(reader io.Reader, slice interface{}) *SliceDeserializer {
	decoder := json.NewDecoder(reader)

	return &SliceDeserializer{
		reader:  reader,
		slice:   slice,
		decoder: decoder,
	}
}

func (d *SliceDeserializer) deserialize() error {
	sliceValue := reflect.ValueOf(d.slice).Elem()

	for {
		if err := d.decodeAndAppendElement(&sliceValue); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

func (d *SliceDeserializer) decodeAndAppendElement(sliceValue *reflect.Value) error {
	element := reflect.New(sliceValue.Type().Elem()).Interface()

	if err := d.decoder.Decode(element); err != nil {
		return err
	}

	sliceValue.Set(reflect.Append(*sliceValue, reflect.ValueOf(element).Elem()))
	return nil
}
