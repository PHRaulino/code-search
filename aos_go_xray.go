package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/xray"
)

func main() {
    // Define the profile you want to use
    profile := "your-profile-name"

    // Load the AWS configuration with the specified profile
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("us-west-2"),
        config.WithSharedConfigProfile(profile),
    )
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Create an X-Ray client
    svc := xray.NewFromConfig(cfg)

    // Define the trace ID
    traceID := "1-5f84c7c4-abc123456def789012345678" // Substitua pelo seu trace ID

    // Get the trace
    getTraceSummariesInput := &xray.GetTraceSummariesInput{
        StartTime: aws.Time(time.Now().Add(-24 * time.Hour)), // Ajuste conforme necessário
        EndTime:   aws.Time(time.Now()),                     // Ajuste conforme necessário
        FilterExpression: aws.String(fmt.Sprintf("traceId = \"%s\"", traceID)),
    }

    result, err := svc.GetTraceSummaries(context.TODO(), getTraceSummariesInput)
    if err != nil {
        log.Fatalf("failed to get trace summaries, %v", err)
    }

    for _, summary := range result.TraceSummaries {
        fmt.Printf("Trace ID: %s\n", *summary.Id)
        fmt.Printf("Duration: %f seconds\n", *summary.Duration)
        fmt.Printf("Segments: %v\n", summary.Segments)
        fmt.Println("-----")
    }

    // Get detailed information about the trace
    getTraceInput := &xray.BatchGetTracesInput{
        TraceIds: []string{traceID},
    }

    traceResult, err := svc.BatchGetTraces(context.TODO(), getTraceInput)
    if err != nil {
        log.Fatalf("failed to get trace details, %v", err)
    }

    for _, trace := range traceResult.Traces {
        fmt.Printf("Trace ID: %s\n", *trace.Id)
        for _, segment := range trace.Segments {
            fmt.Printf("Segment Document: %s\n", *segment.Document)
        }
    }
}
