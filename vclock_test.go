package vclock

import (
	"testing"
	"time"
)

func TestSameProcess(t *testing.T) {
	a1 := New("a")
	a2 := New("a")

	a2.Tick()

	if !a1.Before(a2) {
		t.Error("a1 should come before a2", a1.clocks, a2.clocks)
	}
}

func TestConcurrent(t *testing.T) {
	a := New("a")
	b := New("b")

	if !a.Concurrent(b) {
		t.Error("Initial processes should be concurrent")
	}

	a.Tick()
	b.Tick()
	b.Merge(a)

	if a.Concurrent(b) {
		t.Error("b should be ahead of a")
	}

	a.Merge(b)
	b.Tick()

	if !a.Concurrent(b) {
		t.Error("a and b should be concurrent")
	}
}

func TestBefore(t *testing.T) {
	a := New("a")
	b := New("b")

	a.Tick()
	b.Tick()
	b.Merge(a)

	if !a.Before(b) {
		t.Error("a should come before b")
	}

	a.Merge(b)
	a.Tick()

	if !b.Before(a) {
		t.Error("b should come before a")
	}

	b.Tick()

	if a.Before(b) || b.Before(a) {
		t.Error("a and b should be concurrent")
	}
}

func TestMerge(t *testing.T) {
	a := New("a")
	b := New("b")
	c := New("c")

	b.Tick()
	a.Merge(b)

	if a.clocks["a"] != 1 {
		t.Error("Receiver should increment its clock")
	}

	if a.clocks["b"] != 1 {
		t.Error("Merge should add missing clocks")
	}

	a.Tick()
	a.Tick()
	a.Merge(b)

	if a.clocks["a"] != 4 {
		t.Error("Merge should keep the highest value for a clock")
	}

	b.Merge(c)

	c.Tick()
	c.Tick()
	a.Merge(c)
	a.Merge(b)

	if a.clocks["c"] != 2 {
		t.Error("Merge should keep the highest value for a clock")
	}

	if a.clocks["b"] != 2 {
		t.Error("Merge should increment to higher values")
	}
}

type Process struct {
	*VClock
	socket chan *VClock
}

func (p *Process) send(other *Process) {
	other.socket <- p.VClock
}

func NewProcess(id string) *Process {
	p := &Process{New(id), make(chan *VClock)}

	go func() {
		for {
			clock := <-p.socket
			p.Merge(clock)
		}
	}()

	return p
}

func TestProcesses(t *testing.T) {
	Pi := NewProcess("Pi")
	Pj := NewProcess("Pj")

	Pi.Tick()
	Pj.Tick()
	Pi.send(Pj)

	if !Pi.Before(Pj.VClock) {
		t.Error("Pi should be behind")
	}

	Pi.Tick()

	if !Pi.Concurrent(Pj.VClock) {
		t.Error("Pi and Pj should be concurrent")
	}

	Pj.send(Pi)

	time.Sleep(time.Millisecond * 10)

	if !Pj.Before(Pi.VClock) {
		t.Error("Pj should be behind")
	}
}
