package ip

import "testing"

func testIPAddressExtracer(t *testing.T, testIP string, expectedResult string) {
	result, err := ExtractIPAddress(testIP)
	if err != nil {
		if expectedResult != "ERROR" {
			t.Errorf("function returned an error")
		} else if err != nil && expectedResult == "ERROR" {
			return
		}
	}
	if result != expectedResult {
		t.Errorf("ExtractIPAddress(%s) = %s, expected %s", testIP, result, expectedResult)
	}
}

func TestExtractIPAddress(t *testing.T) {
	testIPAddressExtracer(t, "192.168.1.1:22", "192.168.1.1")
	testIPAddressExtracer(t, "192.168.999.888:5432", "ERROR")
	testIPAddressExtracer(t, "0.0.0.0:9999", "0.0.0.0")
	testIPAddressExtracer(t, "255.255.255.255:2", "255.255.255.255")
}
