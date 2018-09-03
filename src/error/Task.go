package error

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

/*func (err ErrNegativeSqrt) String() string {
	return fmt.Sprintf("%v", float64(err))
}*/

func (err ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", err)
}
func Sqrt(val float64) (float64, error) {
	if val < 0 {
		err := ErrNegativeSqrt(val)
		return 0, err
	}
	toReturn, _ := nuton(val/2, val)
	return toReturn, nil
}

func nuton(z, x float64) (float64, float64) {
	newZ := z - (z*z-x)/(2*z)
	if math.Abs(newZ-z) < 0.00001 {
		return newZ, x
	}
	return nuton(newZ, x)
}

func Task() {

	fmt.Println(nuton(2.5, 5))
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
