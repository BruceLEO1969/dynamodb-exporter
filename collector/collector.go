package collector

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "dynamodb"
)

// Exporter collects metrics from aws dynamodb.
type Exporter struct {
	tablesCount *prometheus.Desc
	// limit       *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter() *Exporter {
	return &Exporter{
		tablesCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "tablesCount"),
			"Count of dynamodb tables with the current account and endpoint.",
			[]string{"aws_region"},
			nil,
		),
		// 	limit: prometheus.NewDesc(
		// 		prometheus.BuildFQName(namespace, "", "limit"),
		// 		"Count of dynamodb tables with the current account and endpoint.",
		// 		[]string{"aws_region"},
		// 		nil,
		// 	),
	}
}

// Describe all the metrics.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.tablesCount
	// ch <- e.limit
}

// Collect fetches the statistics from aws dynamodb sdk, and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	awsRegions := []string{"us-east-1", "ap-northeast-2", "us-east-2", "us-west-1", "us-west-2", "eu-central-1", "eu-west-1"}

	wg := sync.WaitGroup{}

	for _, regionName := range awsRegions {
		regionName := regionName
		wg.Add(1)

		svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion(regionName))
		// svcQuotas := servicequotas.New(session.New(), aws.NewConfig().WithRegion(regionName))

		go func() {
			defer wg.Done()

			// serviceCode := "dynamodb"
			// quotasInput := &servicequotas.ListServiceQuotasInput{}
			// quotasInput := &servicequotas.ListAWSDefaultServiceQuotasInput{}

			// quotasInput.ServiceCode = &serviceCode
			// quotasOutput, err := svcQuotas.ListServiceQuotas(quotasInput)
			// quotasOutput, err := svcQuotas.ListAWSDefaultServiceQuotas(quotasInput)

			// if err != nil {
			// 	fmt.Println(err)
			// }
			// fmt.Printf("quotasOutput in %v:", regionName)
			// fmt.Println(quotasOutput)

			// limitIntput := &dynamodb.DescribeLimitsInput{}
			// limitOutput, err := svc.DescribeLimits(limitIntput)
			// if err != nil {
			// 	return
			// }
			// theLimit := limitOutput.AccountMaxReadCapacityUnits
			// fmt.Printf("limitOutput in %v:", regionName)
			// fmt.Println(limitOutput)
			// ch <- prometheus.MustNewConstMetric(e.limit, prometheus.GaugeValue, float64(*theLimit), regionName)

			// create the input configuration instance
			input := &dynamodb.ListTablesInput{}

			allTables := []string{}
			// Get table list size
			for {
				// Get the list of tables
				result, err := svc.ListTables(input)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						case dynamodb.ErrCodeInternalServerError:
							fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
						default:
							fmt.Println(aerr.Error())
						}
					} else {
						fmt.Println(err.Error())
					}
					return
				}

				for _, n := range result.TableNames {
					allTables = append(allTables, *n)
				}

				// assign the last read tablename as the start for our next call to the ListTables function
				input.ExclusiveStartTableName = result.LastEvaluatedTableName

				if result.LastEvaluatedTableName == nil {
					break
				}
			}

			tablesSize := len(allTables)
			fmt.Printf("Count of tables in %v: %v\n", regionName, tablesSize)

			ch <- prometheus.MustNewConstMetric(e.tablesCount, prometheus.GaugeValue, float64(tablesSize), regionName)
		}()
	}
	wg.Wait()
}
