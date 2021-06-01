package json6

func UnmarshalBytes(src []byte, v interface{}) error {
	obj, err := scanFromBytes(src)
	if err != nil {
		return err
	}

	return unmarshal(obj, v)
}

func unmarshal(obj *object, v interface{}) error {
	return nil
}
