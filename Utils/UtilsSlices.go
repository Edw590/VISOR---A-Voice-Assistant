/*******************************************************************************
 * Copyright 2023-2023 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package Utils

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

/*
DelElemSLICES removes an element from a slice by its index.

If the index is out of range (index < 0 || index >= len(slice)), nothing happens.

Credits to https://stackoverflow.com/a/56591107/8228163 (optimized here).

-----------------------------------------------------------

– Params:
  - slice – a pointer to the slice header
  - index – the index of the element to remove

– Returns:
  - true if the element was removed, false otherwise
*/
func DelElemSLICES(slice any, index int) bool {
	var slice_value reflect.Value = reflect.ValueOf(slice).Elem()

	if index < 0 || index >= slice_value.Len() {
		return false
	}

	slice_value.Set(reflect.AppendSlice(slice_value.Slice(0, index), slice_value.Slice(index+1, slice_value.Len())))

	return true
}

/*
AddElemSLICES adds an element to a specific index of a slice, keeping the elements' order.

-----------------------------------------------------------

– Params:
  - slice – a pointer to the slice header
  - element – the element to add
  - index – the index to add the element on, with range [0, len(slice)]

– Returns:
  - nothing
*/
func AddElemSLICES[T any](slice *[]T, element T, index int) {
	var slice_value reflect.Value = reflect.ValueOf(slice).Elem()
	var element_value reflect.Value = reflect.ValueOf(element)
	var result reflect.Value
	if index > 0 {
		result = reflect.AppendSlice(slice_value.Slice(0, index), slice_value.Slice(index-1, slice_value.Len()))
		result.Index(index).Set(element_value)
	} else {
		var element_slice reflect.Value = reflect.MakeSlice(reflect.SliceOf(element_value.Type()), 1, slice_value.Len()+1)
		element_slice.Index(0).Set(element_value)
		result = reflect.AppendSlice(element_slice, slice_value.Slice(0, slice_value.Len()))
	}
	slice_value.Set(result)
}

/*
CopyOuterSLICES copies all the values from an OUTER slice to a new slice internally created, with the length and
capacity of the original.

Note: the below described won't have any effect if the slice to copy has only one dimension - in that case,
don't worry at all as the function will copy all values normally. If the slice has more dimensions, read the
below explanation.

I wrote “Outer“ in caps because of this example:

var example [][]int = [][]int{{1}, {2}, {3}}

This function will copy the values of the outer slice only - which are pointers to the inner slices. If ANY
value of the inner slices gets changed, on the original slice that shall happen too, because both the original
and the copy point to the same inner slices. Only the outer slices differ - so one can add an inner slice to the
copy, and it will not show up on the original, and change values on that new inner slice - as long as the values
of the original inner slices don't change.

-----------------------------------------------------------

– Params:
  - slice – the slice

– Returns:
  - the new slice
*/
func CopyOuterSLICES[T any](slice T) T {
	var slice_value reflect.Value = reflect.ValueOf(slice)
	var new_slice reflect.Value = reflect.MakeSlice(slice_value.Type(), slice_value.Len(), slice_value.Cap())
	reflect.Copy(new_slice, slice_value)

	return new_slice.Interface().(T)
}

/*
CopyFullSLICES copies all the values from slice/array to a provided slice/array with the length and capacity of the
original.

Note 1: both slices/arrays must have the same type (that includes the length of each dimension with arrays).

NOTE 2: this function is slow, according to what someone told me. Don't use unless it's really necessary to copy all
values from multidimensional slices/arrays.

-----------------------------------------------------------

– Params:
  - dest – a pointer to an empty destination slice/array header
  - src – the source slice/array

– Returns:
  - true if the slice was fully copied, false if an error occurred
*/
func CopyFullSLICES[T any](dest *T, src T) bool {
	var buf *bytes.Buffer = new(bytes.Buffer)
	var err error = gob.NewEncoder(buf).Encode(src)
	if nil != err {
		return false
	}
	err = gob.NewDecoder(buf).Decode(dest)
	if nil != err {
		return false
	}

	return true
}

/*
ContainsSLICES checks if a slice contains a specific element.

From https://stackoverflow.com/a/70802740/8228163.

-----------------------------------------------------------

– Params:
  - s – the slice
  - e – the element to check

– Returns:
  - true if the slice contains the element, false otherwise
*/
func ContainsSLICES[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}

	return false
}
