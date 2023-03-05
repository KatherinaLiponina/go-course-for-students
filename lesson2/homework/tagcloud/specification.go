package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	queue []TagStat
	len int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() * TagCloud {
	var t TagCloud
	t.len = 0
	t.queue = []TagStat{}
	return &t
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (t * TagCloud) AddTag(tag string) {
	for i:= 0; i < t.len; i++ {
		if t.queue[i].Tag == tag {
			t.queue[i].OccurrenceCount++
			for i > 0 && t.queue[i].OccurrenceCount > t.queue[i - 1].OccurrenceCount {
				tmp := t.queue[i];
				t.queue[i] = t.queue[i-1]
				t.queue[i-1] = tmp
			}
			return
		}
	}
	t.len++
	t.queue = append(t.queue, TagStat{tag, 1})
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (t * TagCloud) TopN(n int) []TagStat {
	if n >= t.len {
		return t.queue
	}
	result := []TagStat{}
	for i := 0; i < n; i++ {
		result = append(result, t.queue[i])
	}
	return result
}
