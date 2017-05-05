// Copyright 2016 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"github.com/netsec-ethz/scion/go/lib/common"
)

func CalcPadding(length, blkSize int) int {
	spare := length % blkSize
	if spare != 0 {
		return blkSize - spare
	}
	return 0
}

func FillPadding(b common.RawBytes, length, blkSize int) int {
	padding := CalcPadding(length, blkSize)
	total := length + padding
	for i := range b[length:total] {
		b[i] = 0
	}
	return total
}