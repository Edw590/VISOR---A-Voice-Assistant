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

package UtilsSWA

import (
	"math"
)

/*
IsOutlierMATH checks if an element is an outlier of an array.

The functions does so by checking if the elements is inside a range of mean +- X * standard deviation.

-----------------------------------------------------------

– Params:
  - value – the element to check
  - sum – the sum of all the elements
  - sum_squares – the sum of the squares of all the elements
  - n – the number of elements
  - accuracy_parameter – the mentioned X value

– Returns:
  - true if it's an outlier, false otherwise

*/
func IsOutlierMATH(value float64, sum float64, sum_squares float64, n int32, accuracy_parameter float64) bool {
	mean := sum / float64(n)
	std_dev := math.Sqrt(sum_squares/float64(n) - mean*mean)

	return (value < mean - accuracy_parameter * std_dev) || (value > mean + accuracy_parameter * std_dev)
}
