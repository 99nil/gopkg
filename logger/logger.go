// Copyright © 2021 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

// Interface defines a full log interface
type Interface interface {
	InfoInterface
	WarnInterface
	ErrorInterface
	FatalInterface
	PanicInterface
	DebugInterface
}

// FieldInterface defines a full log with field interface
type FieldInterface interface {
	Interface

	WithField(key string, val interface{}) FieldInterface
	WithFields(fields map[string]interface{}) FieldInterface
}

// UniversalInterface defines a generic log interface
type UniversalInterface interface {
	InfoInterface
	WarnInterface
	ErrorInterface
	DebugInterface
}

// UniversalFieldInterface defines a generic log with field interface
type UniversalFieldInterface interface {
	UniversalInterface

	WithField(key string, val interface{}) UniversalFieldInterface
	WithFields(fields map[string]interface{}) UniversalFieldInterface
}

type StdInterface interface {
	PrintInterface
	FatalInterface
	PanicInterface
}

type PrintInterface interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type InfoInterface interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
}

type WarnInterface interface {
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})
}

type ErrorInterface interface {
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
}

type FatalInterface interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
}

type PanicInterface interface {
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type DebugInterface interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
}

func NewEmpty() FieldInterface {
	return NewEmptyWithStd(NewStdEmpty())
}

func NewEmptyWithStd(std StdInterface) FieldInterface {
	return &empty{std: std}
}

type empty struct {
	std StdInterface
}

func (l *empty) WithField(key string, val interface{}) FieldInterface { return l }

func (l *empty) WithFields(fields map[string]interface{}) FieldInterface { return l }

func (l *empty) Info(v ...interface{})                 { l.std.Print(v...) }
func (l *empty) Infof(format string, v ...interface{}) { l.std.Printf(format, v...) }
func (l *empty) Infoln(v ...interface{})               { l.std.Println(v...) }

func (l *empty) Fatal(v ...interface{})                 { l.std.Fatal(v...) }
func (l *empty) Fatalf(format string, v ...interface{}) { l.std.Fatalf(format, v...) }
func (l *empty) Fatalln(v ...interface{})               { l.std.Fatalln(v...) }

func (l *empty) Panic(v ...interface{})                 { l.std.Panic(v...) }
func (l *empty) Panicf(format string, v ...interface{}) { l.std.Panicf(format, v...) }
func (l *empty) Panicln(v ...interface{})               { l.std.Panicln(v...) }

func (l *empty) Warn(v ...interface{})                 {}
func (l *empty) Warnf(format string, v ...interface{}) {}
func (l *empty) Warnln(v ...interface{})               {}

func (l *empty) Error(v ...interface{})                 {}
func (l *empty) Errorf(format string, v ...interface{}) {}
func (l *empty) Errorln(v ...interface{})               {}

func (l *empty) Debug(v ...interface{})                 {}
func (l *empty) Debugf(format string, v ...interface{}) {}
func (l *empty) Debugln(v ...interface{})               {}

func NewStdEmpty() StdInterface {
	return &stdEmpty{}
}

type stdEmpty struct{}

func (l *stdEmpty) Print(v ...interface{})                 {}
func (l *stdEmpty) Printf(format string, v ...interface{}) {}
func (l *stdEmpty) Println(v ...interface{})               {}

func (l *stdEmpty) Fatal(v ...interface{})                 {}
func (l *stdEmpty) Fatalf(format string, v ...interface{}) {}
func (l *stdEmpty) Fatalln(v ...interface{})               {}

func (l *stdEmpty) Panic(v ...interface{})                 {}
func (l *stdEmpty) Panicf(format string, v ...interface{}) {}
func (l *stdEmpty) Panicln(v ...interface{})               {}
