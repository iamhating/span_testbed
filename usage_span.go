package main

type UsageSpan struct {
	Id string
	Name string
	StartTS int64
	EndTS *int64
}

func (s UsageSpan) GetStartTS() int64 {
	return s.StartTS
}

func (s UsageSpan) GetEndTS() *int64 {
	return s.EndTS
}

func (s UsageSpan) SetStartTS(t int64) {
	s.StartTS = t
}

func (s UsageSpan) SetEndTS(tp *int64) {
	s.EndTS = tp
}

func (s UsageSpan) Copy() Span {
	return UsageSpan{
		Id: s.Id,
		Name: s.Name,
		StartTS: s.StartTS,
		EndTS: s.EndTS,
	}
}