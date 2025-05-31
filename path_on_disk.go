package pathlib


func (p PathOnDisk[any]) Exists() bool {
	return p.Info != nil
}

func (p PathOnDisk[any]) IsRegular() bool {
	return p.Info.Mode().IsRegular()
}
