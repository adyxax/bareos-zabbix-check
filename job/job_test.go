package job

import "testing"

func TestString(t *testing.T) {
	j := Job{Name: "name", Timestamp: 10, Success: true}
	if j.String() != "Job { Name: \"name\", Timestamp: \"10\", Success: \"true\" }" {
		t.Errorf("test string error : %s", j.String())
	}
}
