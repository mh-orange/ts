package ts

type Descriptor []byte

func (d Descriptor) Tag() uint8 {
	return uint8(d[0])
}

func (d Descriptor) Length() int {
	return int(d[1])
}

func (d Descriptor) Data() []byte {
	return d[2:d.Length()]
}

func Descriptors(data []byte) []Descriptor {
	descriptors := make([]Descriptor, 0)
	for i := 0; i < len(data); {
		descriptor := Descriptor(data[i:])
		descriptors = append(descriptors, descriptor)
		i += descriptor.Length() + 2
	}
	return descriptors
}
