package main

import (
	"fmt"
	"time"

	"github.com/gaurav2721/notification-service/notification_manager/scheduler"
	"github.com/sirupsen/logrus"
)

// SchedulerExample demonstrates how to use the job scheduler
func SchedulerExample() {
	fmt.Println("=== Job Scheduler Example ===")

	// Create a new scheduler
	sched := scheduler.NewScheduler()
	defer sched.Stop()

	// Schedule a job to run in 2 seconds
	jobID := "test-job-1"
	scheduledTime := time.Now().Add(2 * time.Second)

	fmt.Printf("Scheduling job '%s' to run at %s\n", jobID, scheduledTime.Format("15:04:05"))

	// Define the job function
	job := func() {
		fmt.Printf("🎯 Job '%s' executed at %s\n", jobID, time.Now().Format("15:04:05"))
		logrus.WithField("job_id", jobID).Info("Scheduled job executed successfully")
	}

	// Schedule the job
	err := sched.ScheduleJob(jobID, scheduledTime, job)
	if err != nil {
		fmt.Printf("❌ Failed to schedule job: %v\n", err)
		return
	}

	fmt.Printf("✅ Job '%s' scheduled successfully\n", jobID)

	// Schedule another job to run in 5 seconds
	jobID2 := "test-job-2"
	scheduledTime2 := time.Now().Add(5 * time.Second)

	fmt.Printf("Scheduling job '%s' to run at %s\n", jobID2, scheduledTime2.Format("15:04:05"))

	job2 := func() {
		fmt.Printf("🎯 Job '%s' executed at %s\n", jobID2, time.Now().Format("15:04:05"))
		logrus.WithField("job_id", jobID2).Info("Second scheduled job executed successfully")
	}

	err = sched.ScheduleJob(jobID2, scheduledTime2, job2)
	if err != nil {
		fmt.Printf("❌ Failed to schedule job: %v\n", err)
		return
	}

	fmt.Printf("✅ Job '%s' scheduled successfully\n", jobID2)

	// List all scheduled jobs
	jobs := sched.GetScheduledJobs()
	fmt.Printf("📋 Currently scheduled jobs: %v\n", jobs)

	// Wait for jobs to execute (6 seconds)
	fmt.Println("⏳ Waiting for jobs to execute...")
	time.Sleep(6 * time.Second)

	// Cancel a job (this will be a no-op since jobs have already executed)
	fmt.Printf("🔄 Cancelling job '%s'...\n", jobID)
	err = sched.CancelJob(jobID)
	if err != nil {
		fmt.Printf("❌ Failed to cancel job: %v\n", err)
	} else {
		fmt.Printf("✅ Job '%s' cancelled successfully\n", jobID)
	}

	fmt.Println("=== Scheduler Example Completed ===")
}

// NotificationSchedulerExample demonstrates scheduling notification jobs
func NotificationSchedulerExample() {
	fmt.Println("=== Notification Scheduler Example ===")

	// Create a new scheduler
	sched := scheduler.NewScheduler()
	defer sched.Stop()

	// Simulate scheduling a notification to be sent in 3 seconds
	notificationID := "notif-123"
	scheduledTime := time.Now().Add(3 * time.Second)

	fmt.Printf("📧 Scheduling notification '%s' to be sent at %s\n", notificationID, scheduledTime.Format("15:04:05"))

	// Define the notification job
	notificationJob := func() {
		fmt.Printf("📧 Sending notification '%s' at %s\n", notificationID, time.Now().Format("15:04:05"))
		logrus.WithField("notification_id", notificationID).Info("Scheduled notification sent successfully")

		// In a real implementation, this would:
		// 1. Fetch the notification details from storage
		// 2. Process the template if needed
		// 3. Send via appropriate channel (email, Slack, push, etc.)
		// 4. Update notification status
		// 5. Log the result
	}

	// Schedule the notification
	err := sched.ScheduleJob(notificationID, scheduledTime, notificationJob)
	if err != nil {
		fmt.Printf("❌ Failed to schedule notification: %v\n", err)
		return
	}

	fmt.Printf("✅ Notification '%s' scheduled successfully\n", notificationID)

	// Wait for notification to be sent
	fmt.Println("⏳ Waiting for notification to be sent...")
	time.Sleep(4 * time.Second)

	fmt.Println("=== Notification Scheduler Example Completed ===")
}

// OneTimeJobExample demonstrates one-time job execution with automatic cleanup
func OneTimeJobExample() {
	fmt.Println("=== One-Time Job Example ===")

	// Create a new scheduler
	sched := scheduler.NewScheduler()
	defer sched.Stop()

	// Track job execution count
	executionCount := 0
	jobID := "one-time-job"

	// Schedule a job to run in 3 seconds
	scheduledTime := time.Now().Add(3 * time.Second)
	fmt.Printf("Scheduling one-time job '%s' to run at %s\n", jobID, scheduledTime.Format("15:04:05"))

	// Define the job function
	job := func() {
		executionCount++
		fmt.Printf("🎯 One-time job '%s' executed (count: %d) at %s\n", jobID, executionCount, time.Now().Format("15:04:05"))
	}

	// Schedule the job
	err := sched.ScheduleJob(jobID, scheduledTime, job)
	if err != nil {
		fmt.Printf("❌ Failed to schedule job: %v\n", err)
		return
	}

	fmt.Printf("✅ One-time job '%s' scheduled successfully\n", jobID)

	// Check scheduled jobs before execution
	jobs := sched.GetScheduledJobs()
	fmt.Printf("📋 Scheduled jobs before execution: %v\n", jobs)

	// Wait for job to execute (4 seconds)
	fmt.Println("⏳ Waiting for job to execute...")
	time.Sleep(4 * time.Second)

	// Check scheduled jobs after execution
	jobs = sched.GetScheduledJobs()
	fmt.Printf("📋 Scheduled jobs after execution: %v\n", jobs)

	// Wait a bit more to ensure no re-execution
	fmt.Println("⏳ Waiting additional 2 seconds to ensure no re-execution...")
	time.Sleep(2 * time.Second)

	fmt.Printf("📊 Final execution count: %d\n", executionCount)
	fmt.Printf("📋 Final scheduled jobs: %v\n", sched.GetScheduledJobs())

	if executionCount == 1 && len(sched.GetScheduledJobs()) == 0 {
		fmt.Println("✅ SUCCESS: Job executed exactly once and was automatically removed!")
	} else {
		fmt.Println("❌ FAILED: Job did not behave as expected")
	}

	fmt.Println("=== One-Time Job Example Completed ===")
}
