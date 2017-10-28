package helper

//InArray Check if a string exists within a string array
//
//@param val: string
//@param array: []string
//
//@returns ok: Defines whether the value exists
//@returns i: Defines whether the position in which it exists
func InArray(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}
