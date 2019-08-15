/* Copyright 2018 The velas Authors
* This file is part of the velas library.
*
* The velas library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The velas library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
*
* You should have received a copy of the GNU Lesser General Public License
* along with the velas library. If not, see <http://www.gnu.org/licenses/>.
 */

package helpers

var emptyHash [32]byte

// ToHash convert slice to hash
func ToHash(in []byte) [32]byte {
	var hash [32]byte
	copy(hash[:], in[:32])
	return hash
}

// HashIsEmpty check hash is empty
func HashIsEmpty(hash [32]byte) bool {
	return hash == emptyHash
}
