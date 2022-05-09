package main

import "fmt"

func toInt64P(i int64) *int64 {
	return &i
}

func main() {
	usage := UsageSpan{
		Id: "test-usage-span",
		Name: "test-usage-span",
		StartTS: int64(100),
		EndTS: toInt64P(200),
	}
	property := PropertySpan{
		Id: "property-of-usage",
		Name: "property-of-usage",
		StartTS: int64(150),
		EndTS: toInt64P(250),
	}

	fmt.Printf("%#v\n", usage)
	fmt.Printf("SpanContainsTime(usage, 150); want: true; got: %t\n", SpanContainsTime(usage, 150))
	fmt.Printf("SpanContainsTime(usage, 50); want: false; got: %t\n", SpanContainsTime(usage, 50))

	fmt.Printf("%#v\n", property)
	fmt.Printf("SpanContainsTime(property, 150); want: true; got: %t\n", SpanContainsTime(property, 150))
	fmt.Printf("SpanContainsTime(property, 50); want: false; got: %t\n", SpanContainsTime(property, 50))


	fmt.Printf("\ncalculating span overlaps\n")
	// I do not like the fact that I have to pass this as []Span, a list of interfaces
	usages := []Span{usage}
	properties := []Span{property}
	overlaps := GetSpanOverlaps(usages, properties)

	fmt.Printf("%#v\n", overlaps)
	fmt.Printf("start: %d, end: %d\n", overlaps[0].GetStartTS(), *overlaps[0].GetEndTS())

	fmt.Printf("\nconverting span overlaps as interfaces to types\n")
	// I do not like how I have to convert interfaces back out to types
	var typedOverlaps []UsageSpan
	for _, overlap := range overlaps {
		if typedOverlap, ok := overlap.(UsageSpan); ok {
			typedOverlaps = append(typedOverlaps, typedOverlap)
		}
		
	}

	for _, typedOverlaps := range typedOverlaps {
		fmt.Printf("\nprinting typed overlap element\n")
		fmt.Printf("%T\n", typedOverlaps)
		fmt.Printf("%#v\n", typedOverlaps)
		fmt.Println(typedOverlaps.Id)
	}	
}