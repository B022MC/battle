package plaza

import "sync"

type ForbidTagStack struct {
	mut  sync.Mutex
	tags []*ForbidMemberTag
}

func (f *ForbidTagStack) Push(tag *ForbidMemberTag) {
	f.mut.Lock()
	f.tags = append(f.tags, tag)
	f.mut.Unlock()
}

func (f *ForbidTagStack) Pop() *ForbidMemberTag {
	f.mut.Lock()
	defer f.mut.Unlock()
	l := len(f.tags)

	if l == 0 {
		return nil
	}

	tag := f.tags[l-1]
	f.tags = f.tags[:l-1]

	return tag
}

func (f *ForbidTagStack) Clear() {
	f.mut.Lock()
	f.tags = []*ForbidMemberTag{}
	f.mut.Unlock()
}
