// File provided by the K Framework Go backend. Timestamp: 2019-07-15 13:11:08.386

package ieletestinginterpreter

import (
	m "github.com/ElrondNetwork/elrond-vm/iele/elrond/standalone/iele-testing-kompiled/ieletestingmodel"
)

// GetNrSteps yields how many steps were executed until now from the start of the last execution
func (i *Interpreter) GetNrSteps() int {
    return i.currentStep
}

// GetLastState yields the current (last) state of the interpreter
func (i *Interpreter) GetState() m.KReference {
     return i.state
}
