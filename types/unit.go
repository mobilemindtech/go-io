package types

type Unit struct {
}

type UnitPtr = *Unit

func OfUnit() *Unit {
	return &Unit{}
}
