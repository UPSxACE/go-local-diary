package testhelper;

type TestHelper[InputType any, ExpectedOutputType any] struct {
	testCaseInputs []InputType;
	testCaseOutputs []ExpectedOutputType;
}

func CreateTestHelper[InputType any, ExpectedOutputType any]() *TestHelper[InputType, ExpectedOutputType]{
	return &TestHelper[InputType, ExpectedOutputType]{}
}

func (thelper *TestHelper[InputType, ExpectedOutputType]) AddTestcase(input InputType, output ExpectedOutputType){
	thelper.testCaseInputs = append(thelper.testCaseInputs, input)
	thelper.testCaseOutputs = append(thelper.testCaseOutputs, output)
}

func (thelper *TestHelper[InputType, ExpectedOutputType]) GetInputs() []InputType{
	return thelper.testCaseInputs
}

func (thelper *TestHelper[InputType, ExpectedOutputType]) GetOutput(index int) ExpectedOutputType{
	return thelper.testCaseOutputs[index]
}
