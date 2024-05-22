package other

type Computation struct {
	f func()
}

func NewComputation(f func()) *Computation {
	return &Computation{f}
}

func (this *Computation) Run() {
	this.f()
}

type Test struct {
	f    func() bool
	test bool
}

func NewTest(f func() bool) *Test {
	return &Test{f: f}
}

func (this *Test) IsTrue() bool {
	return this.test
}

func (this *Test) Run() *Test {
	return &Test{f: this.f, test: this.f()}
}
