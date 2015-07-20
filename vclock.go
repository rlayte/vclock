package vclock

type VClock struct {
	Id     string
	clocks map[string]int
}

func (vc *VClock) Before(other *VClock) bool {
	for id, clock := range vc.clocks {
		if other.clocks[id] < clock {
			return false
		}
	}

	return true
}

func (vc *VClock) Concurrent(other *VClock) bool {
	selfAhead := 0
	otherAhead := 0

	for id, clock := range vc.clocks {
		if other.clocks[id] <= clock {
			selfAhead++
		}
	}

	for id, clock := range other.clocks {
		if vc.clocks[id] <= clock {
			otherAhead++
		}
	}

	return (selfAhead > 0) && (otherAhead > 0)
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
	return &VClock{
		Id:     id,
		clocks: map[string]int{},
	}
}
