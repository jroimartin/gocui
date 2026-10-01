// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestViewWriteCRLF(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := v.Write([]byte("foo\r\nbar\r\n")); err != nil {
		t.Fatal(err)
	}

	got := v.Buffer()
	want := "foo\nbar\n\n"
	if got != want {
		t.Errorf("Buffer() = %q, want %q", got, want)
	}

	gotLines := v.BufferLines()
	wantLines := []string{"foo", "bar", ""}
	if !reflect.DeepEqual(gotLines, wantLines) {
		t.Errorf("BufferLines() = %#v, want %#v", gotLines, wantLines)
	}
}

func TestViewWriteCROverwrite(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := fmt.Fprint(v, "Item 1\nItem 2\n\rWill be deleted\rItem 4\nItem 5"); err != nil {
		t.Fatal(err)
	}

	gotLines := v.BufferLines()
	wantLines := []string{"Item 1", "Item 2", "Item 4", "Item 5"}
	if !reflect.DeepEqual(gotLines, wantLines) {
		t.Errorf("BufferLines() = %#v, want %#v", gotLines, wantLines)
	}
}

func TestViewWriteCRAcrossChunks(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := v.Write([]byte("line 1\r")); err != nil {
		t.Fatal(err)
	}
	if _, err := v.Write([]byte("\nline 2\r\n")); err != nil {
		t.Fatal(err)
	}

	gotLines := v.BufferLines()
	wantLines := []string{"line 1", "line 2", ""}
	if !reflect.DeepEqual(gotLines, wantLines) {
		t.Errorf("BufferLines() = %#v, want %#v", gotLines, wantLines)
	}
}

func TestViewWriteTrailingCR(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := v.Write([]byte("hello\r")); err != nil {
		t.Fatal(err)
	}

	gotLines := v.BufferLines()
	wantLines := []string{"hello"}
	if !reflect.DeepEqual(gotLines, wantLines) {
		t.Errorf("BufferLines() = %#v, want %#v", gotLines, wantLines)
	}
}

func TestViewClearResetsCR(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := v.Write([]byte("hello\r")); err != nil {
		t.Fatal(err)
	}
	v.Clear()
	if _, err := v.Write([]byte("world\n")); err != nil {
		t.Fatal(err)
	}

	gotLines := v.BufferLines()
	wantLines := []string{"world", ""}
	if !reflect.DeepEqual(gotLines, wantLines) {
		t.Errorf("BufferLines() = %#v, want %#v", gotLines, wantLines)
	}
}

func TestViewReadAfterCRLF(t *testing.T) {
	v := &View{}
	v.ei = newEscapeInterpreter(OutputNormal)

	if _, err := v.Write([]byte("first\r\nsecond\r\n")); err != nil {
		t.Fatal(err)
	}

	buf, err := io.ReadAll(v)
	if err != nil {
		t.Fatal(err)
	}

	got := string(buf)
	want := "first\nsecond\n\n"
	if got != want {
		t.Errorf("ReadAll() = %q, want %q", got, want)
	}
}
