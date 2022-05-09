package main

type PropertySpan struct {
	Id string
	Name string
	StartTS int64
	EndTS *int64
}

func (s PropertySpan) GetStartTS() int64 {
	return s.StartTS
}

func (s PropertySpan) GetEndTS() *int64 {
	return s.EndTS
}

func (s PropertySpan) SetStartTS(t int64) {
	s.StartTS = t
}

func (s PropertySpan) SetEndTS(tp *int64) {
	s.EndTS = tp
}

func (s PropertySpan) Copy() Span {
	return PropertySpan{
		Id: s.Id,
		Name: s.Name,
		StartTS: s.StartTS,
		EndTS: s.EndTS,
	}
}