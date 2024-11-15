package helper

func Int64ToIntPtr(i *int64) *int {
	if i == nil {
		return nil
	}
	val := int(*i)
	return &val
}
