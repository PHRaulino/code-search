package cloudwatchmetrics

import (
    "context"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

// GetAverageDuration calcula a duração média de execução da função Lambda no período especificado
func GetAverageDuration(svc *cloudwatch.Client, functionName string, startTime, endTime time.Time) (float64, error) {
    // Obter as métricas
    result, err := svc.GetMetricStatistics(context.TODO(), &cloudwatch.GetMetricStatisticsInput{
        Namespace:  aws.String("AWS/Lambda"),
        MetricName: aws.String("Duration"),
        Dimensions: []types.Dimension{
            {
                Name:  aws.String("FunctionName"),
                Value: aws.String(functionName),
            },
        },
        StartTime: aws.Time(startTime),
        EndTime:   aws.Time(endTime),
        Period:    aws.Int32(86400),
        Statistics: []types.Statistic{
            types.StatisticAverage,
        },
    })
    if err != nil {
        return 0, err
    }

    // Calcular a média geral
    var totalDuration float64
    for _, point := range result.Datapoints {
        totalDuration += *point.Average
    }
    averageDuration := totalDuration / float64(len(result.Datapoints))

    return averageDuration, nil
}
