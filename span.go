package main

// A span is anything that has start and end timestamps, where the start timestamp is never
//	greater than the end timestamp.
// A copy constructor is required in order to carry over metadata defined by the type
//	satisfying the interface. This is necessary when performing operations between two
//	spans that may result in at least one new span.
type Span interface {
	GetStartTS() int64
	GetEndTS() *int64
	SetStartTS(int64)
	SetEndTS(*int64)
	Copy() Span
}

// Handy methods for getting min/max timestamps that handle open spans

func MinStartTS(start1, start2 int64) int64 {
	if start1 < start2 {
		return start1
	}
	return start2
}

func MaxStartTS(start1, start2 int64) int64 {
	if start1 > start2 {
		return start1
	}
	return start2
}

func MinEndTS(end1, end2 *int64) *int64 {
	if end1 == nil && end2 == nil {
		return nil
	}
	// To make prototyping easier, allocate new memory to avoid pointer reference subtleties
	zero := int64(0)
	endCopy := &zero
	if end1 == nil {
		*endCopy = *end2
	} else if end2 == nil {
		*endCopy = *end1
	} else if *end1 <= *end2  {
		*endCopy = *end1
	} else {
		*endCopy = *end2
	}
	return endCopy
}

func MaxEndTS(end1, end2 *int64) *int64 {
	if end1 == nil || end2 == nil {
		return nil
	}

	// To make prototyping easier, allocate new memory to avoid pointer reference subtleties
	zero := int64(0)
	endCopy := &zero
	if *end1 >= *end2  {
		*endCopy = *end1
	} else {
		*endCopy = *end2
	}
	return endCopy
}

// Properties of a span

func IsValid(span Span) bool {
	return span.GetEndTS() == nil || *span.GetEndTS() >= span.GetStartTS()
}

func IsOpen(span Span) bool {
	return span.GetEndTS() == nil
}

func IsClosed(span Span) bool {
	return span.GetEndTS() != nil
}

func IsZeroLength(span Span) bool {
	return span.GetEndTS() != nil && *span.GetEndTS() == span.GetStartTS()
}

func Length(span Span) *int64 {
	if span.GetEndTS() == nil {
		return nil
	}
	boundedLength := *span.GetEndTS() - span.GetStartTS()
	return &boundedLength
}

// Relationships between a span and a timestamp

func SpanIsBeforeTime(span Span, ts int64) bool {
	return span.GetEndTS() != nil && *span.GetEndTS() > ts
}

func SpanIsAfterTime(span Span, ts int64) bool {
	return ts < span.GetStartTS()
}

func SpanContainsTime(span Span, ts int64) bool {
	return ts >= span.GetStartTS() && (span.GetEndTS() == nil || *span.GetEndTS() >= ts)
}

// Relationships between two spans

func SpansEqual(baseSpan, applySpan Span) bool {
	return baseSpan.GetStartTS() == applySpan.GetStartTS() &&
		((baseSpan.GetEndTS() == nil && applySpan.GetEndTS() == nil) ||
		(baseSpan.GetEndTS() != nil && applySpan.GetEndTS() != nil && *baseSpan.GetEndTS() == *applySpan.GetEndTS()))
}

// If the endpoints are the same, we consider that there is no overlap
func SpanLeftOf(baseSpan, applySpan Span) bool {
	return baseSpan.GetEndTS() != nil && *baseSpan.GetEndTS() <= applySpan.GetStartTS()
}

// If the endpoints are the same, we consider that there is no overlap
func SpanRightOf(baseSpan, applySpan Span) bool {
	return applySpan.GetEndTS() != nil && *applySpan.GetEndTS() <= baseSpan.GetStartTS()
}

func HasOverlap(baseSpan, applySpan Span) bool {
	return !(SpanLeftOf(baseSpan, applySpan) || SpanRightOf(baseSpan, applySpan))
}

// Two spans are adjacent if one span ends at the same time that the other starts
func Adjacent(baseSpan, applySpan Span) bool {
	// By our definition of adjacent spans, we want to make sure that later span
	//	starts the same time that the other span ends.
	// So, 1) the start time we want is the start time of the later span;
	//	2) the end time of the earlier span must be finite
	//	3) the start time of the later span must equal the end time of the earlier span
	laterStart := MaxStartTS(baseSpan.GetStartTS(), applySpan.GetStartTS())
	earlyEnd := MinEndTS(baseSpan.GetEndTS() , applySpan.GetEndTS())
	return earlyEnd != nil && laterStart == *earlyEnd
}

func StartsBefore(baseSpan, applySpan Span) bool {
	return baseSpan.GetStartTS() < applySpan.GetStartTS()
}

func StartsAtSameTime(baseSpan, applySpan Span) bool {
	return baseSpan.GetStartTS() == applySpan.GetStartTS()
}

func StartsAfter(baseSpan, applySpan Span) bool {
	return baseSpan.GetStartTS() > applySpan.GetStartTS()
}

func EndsBefore(baseSpan, applySpan Span) bool {
	return baseSpan.GetEndTS() != nil &&
		(applySpan.GetEndTS() == nil || *baseSpan.GetEndTS() < *applySpan.GetEndTS())
}

func EndsAtSameTime(baseSpan, applySpan Span) bool {
	return (baseSpan.GetEndTS() == nil && applySpan.GetEndTS() == nil) ||
		(baseSpan.GetEndTS() != nil && applySpan.GetEndTS() != nil && *baseSpan.GetEndTS() == *applySpan.GetEndTS())
}

func EndsAfter(baseSpan, applySpan Span) bool {
	return applySpan.GetEndTS() != nil &&
		(baseSpan.GetEndTS() == nil || *baseSpan.GetEndTS() > *applySpan.GetEndTS())
}

// Returns a span representing the overlap between two spans
func SpanOverlap(baseSpan, applySpan Span) Span {
	var overlap Span
	if HasOverlap(baseSpan, applySpan) {
		overlap = baseSpan.Copy()
		overlap.SetStartTS(MaxStartTS(baseSpan.GetStartTS(), applySpan.GetStartTS()))
		overlap.SetEndTS(MinEndTS(baseSpan.GetEndTS(), applySpan.GetEndTS()))
	}
	return overlap
}

// Returns a span representing the union of two spans, if spans overlap or are adjacent
// Otherwise, spans cannot be merged and return nil
func SpanMerge(baseSpan, applySpan Span) Span {
	var merge Span
	if HasOverlap(baseSpan, applySpan) || Adjacent(baseSpan, applySpan) {
		merge := baseSpan.Copy()
		merge.SetStartTS(MinStartTS(baseSpan.GetStartTS(), applySpan.GetStartTS()))
		merge.SetEndTS(MaxEndTS(baseSpan.GetEndTS(), applySpan.GetEndTS()))
	}
	return merge
}

// This method diff's out the portion of the base span that overlaps with any
//	portion of the apply span..
// Span difference is interesting since it is possible to have 0, 1 or 2 spans
//	returned from performing span difference.
// If there is overlap, hasOverlap is true. leftDiff and rightDiff will be the
//	portions of the spans that lie to the left and right of the overlap.
//	leftDiff or rightDiff or both may be nil if no such poritions remain.
// If there is no overlap, hasOverlap is false. Both leftDiff and rightDiff
//	will be nil
func SpanDiff(baseSpan, applySpan Span) (leftDiff, rightDiff Span, hasOverlap bool) {
	overlap := SpanOverlap(baseSpan, applySpan)
	if overlap.GetStartTS() > 0 && Length(overlap) == nil {
		// If there is no overlap, there is no span left when taking the difference
		return
	}

	// Set flag to distinguish between the case when leftDiff and rightDiff are both
	//	nil when there is no overlap from when there is no remaining difference after
	//	excluding the overlap
	hasOverlap = true
	
	// Check if there is any portion of the base span that precedes the overlap.
	// If so, this is the left portion of the remaining base span
	if StartsBefore(baseSpan, overlap) {
		leftDiff := baseSpan.Copy()
		leftDiffEnd := overlap.GetStartTS()
		leftDiff.SetEndTS(&leftDiffEnd)
	}
	
	// Check if there is any portion of the base span that succeeds the overlap.
	// If so, this is the right porition of the remaining base span
	if EndsBefore(overlap, baseSpan) {
		rightDiff := baseSpan.Copy()
		rightDiffStart := overlap.GetEndTS()
		rightDiff.SetStartTS(*rightDiffStart)
	}

	return
}

func GetSpanOverlaps(baseSpans, applySpans []Span) (overlaps []Span) {
	baseSpanIdx := 0
	applySpanIdx := 0
	for baseSpanIdx < len(baseSpans) && applySpanIdx < len(applySpans) {
		baseSpan := baseSpans[baseSpanIdx]
		applySpan := applySpans[applySpanIdx]

		if SpanLeftOf(baseSpan, applySpan) {
			baseSpanIdx++
		} else if SpanRightOf(baseSpan, applySpan) {
			applySpanIdx++
		} else {
			overlap := SpanOverlap(baseSpan, applySpan)
			if overlap != nil {
				overlaps = append(overlaps, overlap)
			}
			if EndsAtSameTime(baseSpan, overlap) {
				baseSpanIdx++
			} else {
				applySpanIdx++
			}
		}
	}

	return
}

// You now need to assert that the type of the overlap span is the same as the
//	constrained type T, used to represent the type of the output list
// This can be done at any point where we call a function that returns any type
//	satisfying Span, either in GenericSpanOverlap or GenericGetSpanOverlaps
// There should be logic to ensure that the correct type is returned. Otherwise,
//	the code panics. We ignore this check for now quick demonstration purposes
	
func GenericSpanOverlap[T, U Span](baseSpan T, applySpan U) T {
	var overlap T
	if HasOverlap(baseSpan, applySpan) {
		overlap = baseSpan.Copy().(T) // One place for type casting
		overlap.SetStartTS(MaxStartTS(baseSpan.GetStartTS(), applySpan.GetStartTS()))
		overlap.SetEndTS(MinEndTS(baseSpan.GetEndTS(), applySpan.GetEndTS()))
	}
	return overlap
}

func GenericGetSpanOverlaps[T, U Span](baseSpans []T, applySpans []U) (overlaps []T) {
	baseSpanIdx := 0
	applySpanIdx := 0
	for baseSpanIdx < len(baseSpans) && applySpanIdx < len(applySpans) {
		baseSpan := baseSpans[baseSpanIdx]
		applySpan := applySpans[applySpanIdx]

		if SpanLeftOf(baseSpan, applySpan) {
			baseSpanIdx++
		} else if SpanRightOf(baseSpan, applySpan) {
			applySpanIdx++
		} else {
			overlap := GenericSpanOverlap(baseSpan, applySpan) // An alternative for type casting
			overlaps = append(overlaps, overlap)
			if EndsAtSameTime(baseSpan, overlap) {
				baseSpanIdx++
			} else {
				applySpanIdx++
			}
		}
	}

	return
}
