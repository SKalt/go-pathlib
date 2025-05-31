package pathlib


func (p *PathOnDisk[P]) Parent() Dir {
	return PathStr(p.original).Parent()
}

func (p *PathOnDisk[P]) Path() P {
	return p.original
}
