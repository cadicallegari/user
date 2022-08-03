package xdatabase

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
)

type prevCursorer interface {
	PrevCursor()
}

type nextCursorer interface {
	NextCursor()
}

type Cursor struct {
	// PrimaryKey is used internally to build some of the query parameters. Should not be empty
	PrimaryKey string

	PerPage uint64
	// Page is used to do pagination with limit and offset
	Page uint64

	// SearchBefore is used to do backware with a cursor
	SearchBefore string
	// SearchAfter is used to paginate forward with a cursor
	SearchAfter string

	// Total is a pointer which keep the total of the records
	// it's a pointer because it's ok to load the total of available
	// only once
	Total *uint64

	// Search is any arbitrary data to be searched
	Search string
	// Embed it's a list of related fields to be embed in the result
	Embed []string
}

func (c *Cursor) PrevCursor() {
	if c.Page > 0 {
		c.Page--
	}
}

func (c *Cursor) NextCursor() {
	c.Page++
}

// NextCursor make a copy of c to call NextCursor() and return
// a pointer to a string which is the enconded cursor.
func NextCursor(c nextCursorer) *string {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(c); err != nil {
		panic(fmt.Sprint("unable to encode cursor:", err))
	}

	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	c2 := reflect.New(v.Type()).Interface().(nextCursorer)
	err := gob.NewDecoder(buf).Decode(c2)
	if err != nil || c2 == nil {
		return nil
	}

	c2.NextCursor()
	encc := EncodeCursor(c2)
	if encc == "" {
		return nil
	}
	return &encc
}

// PrevCursor make a copy of c to call PrevCursor() and return
// a pointer to a string which is the enconded cursor.
func PrevCursor(c prevCursorer) *string {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(c); err != nil {
		panic(fmt.Sprint("unable to encode cursor:", err))
	}

	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	c2 := reflect.New(v.Type()).Interface().(prevCursorer)
	err := gob.NewDecoder(buf).Decode(c2)
	if err != nil || c2 == nil {
		return nil
	}

	c2.PrevCursor()
	encc := EncodeCursor(c2)
	if encc == "" {
		return nil
	}
	return &encc
}

func DecodeCursor(encc string, c interface{}) error {
	data, err := base64.URLEncoding.DecodeString(encc)
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(c)
}

func EncodeCursor(c interface{}) string {
	buf := new(bytes.Buffer)
	gob.NewEncoder(buf).Encode(c)
	return base64.URLEncoding.EncodeToString(buf.Bytes())
}
