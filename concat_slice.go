package wordlab

/*
Concatenation functions for slices
*/

// ConcatStringSlice concatenates two slices of strings.
// Takes two arguments that must be slices of type String, order of args is not important.
// Returns a new slice which is the result of copying the 2 slices passed in a args.
//
//  slice_one := []string{"this, that, other"}
//  slice_two := []string{"when, where, why"}
//  ConcatStringSlice(slice_one,slice_two)
func ConcatStringSlice(slice1, slice2 []string) []string {
	new_slice := make([]string, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}

// ConcatFloat32Slice concatenates two slices of float32.
// See documentation for ConcatStringSlice for more detail
func ConcatFloat32Slice(slice1, slice2 []float32) []float32 {
	new_slice := make([]float32, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}

// ConcatByteSlice concatenates two slices of bytes.
// See documentation for ConcatStringSlice for more detail
func ConcatByteSlice(slice1, slice2 []byte) []byte {
	new_slice := make([]byte, len(slice1)+len(slice2))
	copy(new_slice, slice1)
	copy(new_slice[len(slice1):], slice2)
	return new_slice
}
