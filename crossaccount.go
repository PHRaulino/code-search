package main

import (
    "context"
    "fmt"
    "time"

    "cloudwatchmetrics"
    "lambdafunctions"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/cloudwatch"
    "github.com/aws/aws-sdk-go-v2/service/lambda"
)

func main() {
    // Configurações para as duas contas
    profiles := []string{"profile1", "profile2"}
    startTime := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
    endTime := time.Date(2023, 5, 31, 23, 59, 59, 0, time.UTC)

    // Iterar sobre os perfis e obter as métricas
    for _, profile := range profiles {
        cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profile), config.WithRegion("us-east-1"))
        if err != nil {
            panic("unable to load SDK config, " + err.Error())
        }

        // Criar os clientes Lambda e CloudWatch
        lambdaSvc := lambda.NewFromConfig(cfg)
        cloudwatchSvc := cloudwatch.NewFromConfig(cfg)

        // Listar funções Lambda
        functions, err := lambdafunctions.ListFunctions(lambdaSvc)
        if err != nil {
            panic("failed to list functions, " + err.Error())
        }

        fmt.Printf("Métricas para perfil %s:\n", profile)
        for _, functionName := range functions {
            averageDuration, err := cloudwatchmetrics.GetAverageDuration(cloudwatchSvc, functionName, startTime, endTime)
            if err != nil {
                fmt.Printf("Erro ao obter métricas para a função %s: %s\n", functionName, err)
                continue
            }
            fmt.Printf("Função: %s, Tempo médio de execução: %.2f ms\n", functionName, averageDuration)
        }
    }
}
