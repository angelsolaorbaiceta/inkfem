package load

import "fmt"

type LoadTerm string
const (
    FX = LoadTerm("fx")
    FY = LoadTerm("fy")
    MZ = LoadTerm("mz")
)

func IsValidTerm(term LoadTerm) bool {
    return (term == FX) || (term == FY) || (term == MZ)
}

func EnsureValidTerm(term LoadTerm) {
    if !IsValidTerm(term) {
        panic(fmt.Sprintf("Invalid load term: '%s'", term))
    }
}
