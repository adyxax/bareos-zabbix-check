package job

// KeepOldestOnly filters a job list and keeps only the most recent entry for a job name
func KeepOldestOnly(jobs []Job) (results []Job) {
outerLoop:
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		for j := 0; j < len(results); j++ {
			result := results[j]
			if result.Name == job.Name {
				continue outerLoop
			}
		}
		for j := i + 1; j < len(jobs); j++ {
			sec := jobs[j]
			if sec.Name == job.Name && sec.Timestamp > job.Timestamp {
				job = sec
			}
		}
		results = append(results, job)
	}
	return
}

// KeepSuccessOnly returns only the successful jobs from a job list (suiatble to write a new spool file)
func KeepSuccessOnly(jobs []Job) (result []Job) {
	result = make([]Job, 0)
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		if job.Success {
			result = append(result, job)
		}
	}
	return
}
