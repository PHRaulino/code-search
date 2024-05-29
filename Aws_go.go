package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func main() {
    // Load the AWS configuration
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Create a CloudWatch Logs client
    svc := cloudwatchlogs.NewFromConfig(cfg)

    // Define the log group name
    logGroupName := "/aws/lambda/LambdaLogin"

    // Define the query string
    queryString := `fields @timestamp, @message
    | filter @message like /Lambda(Login|Logout|OpendID)/
    | sort @timestamp desc
    | limit 20`

    startTimeStr := "2023-05-01T00:00:00Z"
    endTimeStr := "2023-05-02T00:00:00Z"
    layout := time.RFC3339
    
    startTime, err := time.Parse(layout, startTimeStr)
    if err != nil {
        log.Fatalf("failed to parse start time, %v", err)
    }
    
    endTime, err := time.Parse(layout, endTimeStr)
    if err != nil {
        log.Fatalf("failed to parse end time, %v", err)
    }
    
    // Convert to Unix time (in milliseconds)
    startTimeUnix := startTime.Unix() * 1000
    endTimeUnix := endTime.Unix() * 1000

    // Start the query
    startQueryInput := &cloudwatchlogs.StartQueryInput{
        LogGroupName: aws.String(logGroupName),
        StartTime:    aws.Int64(startTime),
        EndTime:      aws.Int64(endTime),
        QueryString:  aws.String(queryString),
    }

    startQueryOutput, err := svc.StartQuery(context.TODO(), startQueryInput)
    if err != nil {
        log.Fatalf("failed to start query, %v", err)
    }

    queryId := startQueryOutput.QueryId

    // Poll for the query results
    for {
        getQueryResultsInput := &cloudwatchlogs.GetQueryResultsInput{
            QueryId: queryId,
        }

        getQueryResultsOutput, err := svc.GetQueryResults(context.TODO(), getQueryResultsInput)
        if err != nil {
            log.Fatalf("failed to get query results, %v", err)
        }

        if *getQueryResultsOutput.Status == "Complete" {
            for _, result := range getQueryResultsOutput.Results {
                for _, field := range result {
                    fmt.Printf("%s: %s\n", aws.ToString(field.Field), aws.ToString(field.Value))
                }
            }
            break
        }

        // Wait for a while before polling again
        time.Sleep(1 * time.Second)
    }
}
