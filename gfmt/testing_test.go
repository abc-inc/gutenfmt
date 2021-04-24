/**
 * Copyright 2021 The gutenfmt authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gfmt_test

import (
	"math"
	"strings"
)

type User struct {
	Name     string `json:"username"`
	Mail     string `json:"email"`
	Password string `json:"-"`
}

func NewUser(fName, lName string) *User {
	return &User{
		fName + " " + lName,
		strings.ToLower(fName + "." + lName + "@local"),
		"",
	}
}

func (u User) String() string {
	return u.Name + " <" + u.Mail + ">"
}

type Team struct {
	name    string
	members []User
}

func NewTeam(name string, users ...User) *Team {
	return &Team{name, users}
}

func (t Team) Name() string {
	return strings.ToUpper(t.name)
}

func (t Team) Members() []User {
	return t.members
}

type JSONTypes struct {
	DefName        string `json:""`
	Skip           string `json:"-"`
	OmitEmpty      string `json:",omitempty"`
	CustOmitEmpty  string `json:"custom,omitempty"`
	EmptyOmitEmpty string `json:",omitempty"`
	Bool           bool
	Int            int
	Int8           int8
	Int16          int16
	Int32          int32
	Int64          int64
	Uint           uint
	Uint8          uint8
	Uint16         uint16
	Uint32         uint32
	Uint64         uint64
	Uintptr        uintptr
	Float32        float32
	Float64        float64
	Array          [2]string
	Interface      interface{}
	Map            map[string]int
	Ptr            *User
	Slice          []string
	String         string
	Struct         User
	StructSlice    []User
}

type AllTypes struct {
	DefName        string `json:""`
	Skip           string `json:"-"`
	OmitEmpty      string `json:",omitempty"`
	CustOmitEmpty  string `json:"custom,omitempty"`
	EmptyOmitEmpty string `json:",omitempty"`
	Bool           bool
	Int            int
	Int8           int8
	Int16          int16
	Int32          int32
	Int64          int64
	Uint           uint
	Uint8          uint8
	Uint16         uint16
	Uint32         uint32
	Uint64         uint64
	Uintptr        uintptr
	Float32        float32
	Float64        float64
	Complex64      complex64
	Complex128     complex128
	Array          [2]string
	Chan           chan<- int
	Func           func(string, string) *User
	Interface      interface{}
	Map            map[interface{}]interface{}
	Ptr            *User
	Slice          []string
	String         string
	Struct         User
	StructSlice    []User
	// UnsafePointer
}

var allTypes = AllTypes{
	"DefName", "Skip",
	"OmitEmpty", "CustOmitEmpty", "",
	true,
	-4, -8, -16, -32, -64,
	4, 8, 16, 32, 64, 128,
	math.MaxFloat32, math.MaxFloat64,
	complex64(-2.71i), complex128(-3.14i),
	[2]string{"a", "b"},
	make(chan<- int),
	NewUser,
	"",
	make(map[interface{}]interface{}),
	NewUser("f", "l"),
	[]string{"a", "b"},
	"",
	*NewUser("f", "l"),
	[]User{*NewUser("af", "al"), *NewUser("bf", "bl")},
}

var jsonTypes = JSONTypes{
	"DefName", "Skip",
	"OmitEmpty", "CustOmitEmpty", "",
	true,
	-4, -8, -16, -32, -64,
	4, 8, 16, 32, 64, 128,
	math.MaxFloat32, math.MaxFloat64,
	[2]string{"a", "b"},
	"",
	make(map[string]int),
	NewUser("f", "l"),
	[]string{"a", "b"},
	"",
	*NewUser("f", "l"),
	[]User{*NewUser("af", "al"), *NewUser("bf", "bl")},
}
