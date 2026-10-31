package models

// QueueHealth is the part of the job queue that needs a person.
type QueueHealth struct {
	WaitingTooLong int64 // ready to run for more than five minutes
	FailedToday    int64 // gave up after every retry, in the last day
}
