/*
 *   Copyright (c) 2025 Vecble
 *   All rights reserved.

 *   Permission is hereby granted, free of charge, to any person obtaining a copy
 *   of this software and associated documentation files (the "Software"), to deal
 *   in the Software without restriction, including without limitation the rights
 *   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *   copies of the Software, and to permit persons to whom the Software is
 *   furnished to do so, subject to the following conditions:

 *   The above copyright notice and this permission notice shall be included in all
 *   copies or substantial portions of the Software.

 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *   SOFTWARE.
 */

package storage

// Represents object data type of a entry. It could be various types such as JSON, int, set, hashmap, string, ...etc
type Object struct {
	ObjectType ObjectType
	Value      interface{}
}

type ObjectType uint

const (
	ObjecTypeString ObjectType = 1
	ObjectTypeInt   ObjectType = 2
	ObjectTypeSet   ObjectType = 3
	ObjectTypeArray ObjectType = 4
	ObjectTypeList  ObjectType = 5
)

func (o Object) String() string {
	switch o.ObjectType {
	case ObjectTypeInt:
		return "int"
	case ObjectTypeSet:
		return "set"
	case ObjectTypeArray:
		return "array"
	case ObjectTypeList:
		return "list"
	default:
		return "string"
	}
}
