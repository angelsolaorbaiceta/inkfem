package structure

var (
	NIL_CONSTRAINT Constraint = Constraint{false, false, false}
	DISP_CONSTRAINT Constraint = Constraint{true, true, false}
	FULL_CONSTRAINT Constraint = Constraint{true, true, true}
)

type Constraint struct {
	isDxConstr, isDyConstr, isRzConst bool
}

/* Construction */
func MakeConstraint(isDxConstr, isDyConstr, isRzConst bool) Constraint {
	switch {
	case !isDxConstr && !isDyConstr && !isRzConst:
		return NIL_CONSTRAINT

	case isDxConstr && isDyConstr && !isRzConst:
		return DISP_CONSTRAINT

	case isDxConstr && isDyConstr && isRzConst:
		return FULL_CONSTRAINT

	default:
		return Constraint{isDxConstr, isDyConstr, isRzConst}
	}
}

func MakeNilConstraint() Constraint {
	return NIL_CONSTRAINT
}

func MakeDispConstraint() Constraint {
	return DISP_CONSTRAINT
}

func MakeFullConstraint() Constraint {
	return FULL_CONSTRAINT
}

/* Methods */
func (c Constraint) AllowsRotation() bool {
	return !c.isRzConst
}
