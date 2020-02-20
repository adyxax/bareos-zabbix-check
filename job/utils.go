package job

// KeepOldestOnly filters a job list and keeps only the most recent entry for a job name
func KeepOldestOnly(jobs []Job) []Job {
	tmpMap := make(map[string]Job)
	for _, elt := range jobs {
		prev, exists := tmpMap[elt.Name]
		if !exists || (exists && prev.Timestamp < elt.Timestamp) {
			tmpMap[elt.Name] = elt
		}
	}
	values := make([]Job, 0, len(tmpMap))
	for _, value := range tmpMap {
		values = append(values, value)
	}
	return values
}

// KeepSuccessOnly returns only the successful jobs from a job list (suiatble to write a new spool file)
func KeepSuccessOnly(jobs []Job) (result []Job) {
	result = make([]Job, 0)
	for _, job := range jobs {
		if job.Success {
			result = append(result, job)
		}
	}
	return
}
