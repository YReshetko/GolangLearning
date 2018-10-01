package main

import (
	"learning/array"
	"learning/dictionary"
	"learning/error"
	"learning/functions"
	"learning/gorutin"
	"learning/hello"
	"learning/img"
	"learning/interf"
	"learning/methods"
	"learning/reader"
	"learning/regulars"
	"learning/synchronization"
	"learning/winProc"
)

func main() {
	//helloPack()
	//arrayPack()
	//dictionaryPack()
	//functionPack()
	//methodsPack()
	//interfPack()
	//errorPack()
	//readerPack()
	//readerImage()
	//gorutPack()
	//synchPack()
	//regularsPack()
	winProcPack()
}
func winProcPack() {
	winProc.Experiment1()
}
func regularsPack() {
	regulars.Experiment1()
}
func synchPack() {
	//synchronization.Experiment1()
	synchronization.Task()
}

func gorutPack() {
	//gorutin.Experiment1()
	//gorutin.Experiment2()
	//gorutin.Experiment3()
	//gorutin.Experiment4()
	//gorutin.Experiment5()
	//gorutin.Experiment6()
	gorutin.Task()
}
func readerImage() {
	//img.Experiment1()
	img.Task()
}
func readerPack() {
	//reader.Experiment1()
	//reader.Task()
	reader.Task2()
}
func errorPack() {
	//error.Experiment1()
	error.Task()
}
func interfPack() {
	//interf.Experiment1()
	//interf.Experiment2()
	//interf.Experiment3()
	//interf.Experiment4()
	//interf.Experiment5()
	//interf.Experiment6()
	//interf.Experiment7()
	//interf.Experiment8()
	//interf.Experiment9()
	interf.Task()
}
func methodsPack() {
	//methods.Experiment1()
	//methods.Experiment2()
	//methods.Experiment3()
	methods.Experiment4()
}
func functionPack() {
	//functions.Experiment1()
	//functions.Experiment2()
	//functions.Experiment3()
	functions.Task()
}

func dictionaryPack() {
	//dictionary.Experiment1()
	//dictionary.Experiment2()
	//dictionary.Experiment3()
	//dictionary.Experiment4()
	dictionary.Task()
}
func arrayPack() {
	//array.Experiment1()
	array.Experiment2()
}
func helloPack() {
	hello.Experiment1()
}
