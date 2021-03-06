package main

import (
	"fmt"
	"path"

	interpreter "github.com/ElrondNetwork/elrond-vm/iele/elrond/standalone/iele-testing-kompiled/ieletestinginterpreter"
)

func main() {
	if false {
		interpreter.StartTraceReferenceComparer("ocaml_or3_1.log")
	}

	debugTest("tests/iele-v3/danse/ERC20/transferFrom_AllDistinct-NoOverflow.iele.json")
	// debugTest("tests/iele/danse/factorial/factorial_positive.iele.json")
	//debugTest("tests/iele-v2/danse/forwarder/create.iele.json")
}

// .build/vm/iele-test-vm tests/iele/danse/forwarder/create.iele.json 10000

func debugTest(testFile string) {
	err := runTest(path.Join(ieleTestRoot, testFile), gasModeNormal, false)
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
