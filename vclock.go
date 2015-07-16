package vclock

type VClock struct {
	Id     string
	clocks map[string]int
}

func (vc *VClock) Before(other *VClock) (ret bool) {
	for id, clock := range vc.clocks {
		if other.clocks[id] < clock {
			return
		}
	}

	ret = true
	return
}

func (vc *VClock) Concurrent(other *VClock) (ret bool) {
	for id, clock := range vc.clocks {
		val, ok := other.clocks[id]
		if ok && val < clock {
			return
		}
	}

	ret = true
	return
}

func (vc *VClock) Tick() {
	vc.clocks[vc.Id]++
}

func (vc VClock) Merge(other *VClock) {
	vc.Tick()

	for id, clock := range other.clocks {
		if clock > vc.clocks[id] {
			vc.clocks[id] = clock
		}
	}
}

func New(id string) *VClock {
	clocks := map[string]int{}
	clocks[id] = 0

	return &VClock{
		Id:     id,
		clocks: clocks,
	}
}
