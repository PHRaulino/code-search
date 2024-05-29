package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
    // Create a session with the AWS SDK
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"),
    })
    if err != nil {
        fmt.Println("Error creating session:", err)
        return
    }

    // Create a CloudWatch Logs client
    svc := cloudwatchlogs.New(sess)

    // Define the log group name
    logGroupName := "/aws/lambda/LambdaLogin"

    // List log streams
    logStreamsInput := &cloudwatchlogs.DescribeLogStreamsInput{
        LogGroupName: aws.String(logGroupName),
    }

    result, err := svc.DescribeLogStreams(logStreamsInput)
    if err != nil {
        fmt.Println("Error describing log streams:", err)
        return
    }

    for _, stream := range result.LogStreams {
        fmt.Println("Log Stream Name:", *stream.LogStreamName)

        // Get log events for each log stream
        logEventsInput := &cloudwatchlogs.GetLogEventsInput{
            LogGroupName:  aws.String(logGroupName),
            LogStreamName: stream.LogStreamName,
        }

        logEvents, err := svc.GetLogEvents(logEventsInput)
        if err != nil {
            fmt.Println("Error getting log events:", err)
            continue
        }

        for _, event := range logEvents.Events {
            fmt.Println("Log Event:", *event.Message)
        }
    }
}
