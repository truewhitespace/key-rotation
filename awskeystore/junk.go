package awskeystore

//StringSlice wraps a primitive string slice to provide additional utilities
//TODO: See if there is a library which already implements these
type StringSlice []string

//Contains a search of the slice elements to determine if the provided element exists within the slice.  true if the
//element exists according to the infix `==` operator.
func (l StringSlice) Contains(element string) bool {
	for _, v := range l {
		if v == element {
			return true
		}
	}
	return false
}
