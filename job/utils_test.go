package job

import "testing"

func TestKeepOldestOnly(t *testing.T) {
	t.Run("test empty list", func(t *testing.T) {
		var jobs []Job
		if len(KeepOldestOnly(jobs)) != 0 {
			t.Error("empty list failed")
		}
	})
	t.Run("test functionality", func(t *testing.T) {
		var jobs []Job
		jobs = append(jobs, Job{Name: "a", Timestamp: 10, Success: true})
		jobs = append(jobs, Job{Name: "a", Timestamp: 20, Success: true})
		jobs2 := KeepOldestOnly(jobs)
		if len(jobs2) != 1 || jobs2[0].Timestamp != 20 {
			t.Error("functionality failed")
		}
	})
}

func TestKeepSuccessOnly(t *testing.T) {
	t.Run("test empty list", func(t *testing.T) {
		var jobs []Job
		if len(KeepSuccessOnly(jobs)) != 0 {
			t.Error("empty list failed")
		}
	})
	t.Run("test functionality", func(t *testing.T) {
		var jobs []Job
		jobs = append(jobs, Job{Name: "a", Timestamp: 10, Success: true})
		jobs = append(jobs, Job{Name: "b", Timestamp: 20, Success: false})
		if len(KeepSuccessOnly(jobs)) != 1 || jobs[0].Name != "a" {
			t.Error("functionality failed")
		}
	})
}
