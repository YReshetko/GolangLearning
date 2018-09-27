package regulars

import (
	"fmt"
	"regexp"
)

func Experiment1() {
	//exp1 := regexp.MustCompile("выковырять ([0-9]+) затем .+\n+.+")
	exp1 := regexp.MustCompile("выковырять ([0-9]+) затем.+\n.+числовую ([A-Za-z0-9]+),.+ ([0-9]+)")
	var str string = fmt.Sprintf("Привет Юра\n")
	str += fmt.Sprintf("Это тестовый емэйл где тебе нужно выковырять 123443654 затем получить\n")
	str += fmt.Sprintf("тектовую или числовую 3apic6, а потом номер 34235\n")
	str += fmt.Sprintf("Не забудь про аттачмент\n")
	fmt.Printf("%q\n", str)
	fmt.Printf("%q\n", exp1.FindStringSubmatch(str))

	//exp2 := regexp.MustCompile("[0-9]{1,2}[Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec]{1}[0-9]{4}_.*\\.txt")
	//exp2 := regexp.MustCompile("[0-9]{1,2}(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[0-9]{4}_.*\\.txt")
	exp2 := regexp.MustCompile("([0-9]{1,2})(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)([0-9]{4})_([0-9]{1,2})h([0-9]{1,2})m([0-9]{1,2})s\\.txt")
	str1 := "3Sep2018_19h24m12s.txt"
	fmt.Printf("%q\n", str1)
	fmt.Printf("%q\n", exp2.FindStringSubmatch(str1))
}
