package compute

import "testing"

func TestParseMacOSGPUInformation(t *testing.T) {
	gpus, err := parseMacOSGPUInformation([]byte(`{"SPDisplaysDataType":[{"_name":"Fallback","sppci_model":" Apple M4 "},{"_name":""}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(gpus) != 2 || gpus[0].Model != "Apple M4" || gpus[1].Model != unknownGPUModel {
		t.Fatalf("unexpected GPUs: %#v", gpus)
	}
}
