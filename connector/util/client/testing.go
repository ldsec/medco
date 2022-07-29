package utilclient

// InitializeStringPointer declares a new string pointer, assigns a value and returns the pointer
func InitializeStringPointer(s string) *string {
	ret := new(string)
	*ret = s
	return ret
}

// InitializeBoolPointer declares a new bool pointer, assigns a value and returns the pointer
func InitializeBoolPointer(s bool) *bool {
	ret := new(bool)
	*ret = s
	return ret
}

// InitializeInt64Pointer declares a new int64 pointer, assigns a value and returns the pointer
func InitializeInt64Pointer(s int64) *int64 {
	ret := new(int64)
	*ret = s
	return ret
}
