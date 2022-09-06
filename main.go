package main

import (
	"flag"
    "net/http"
    "github.com/ghodss/yaml"
	"os"
	"io/ioutil"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/common/version"
    log "github.com/sirupsen/logrus"

    "fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
    "time"
)

type Conf struct {
	Aws	        aws_config              `yaml:"AWS"`
    Metrics     map[string]struct {
        Description string              `yaml:"description"`
        Type        string              `yaml:"type"`
        metricDesc  *prometheus.Desc
    }
}

type aws_config struct {
    Period  	int64                   `yaml:"period"`
    Logmode     string                  `yaml:"logmode"`
}

var config Conf
var region string 
var instance string

const (
	collector = "cloud_exporter"
)

func main() {
    // Get OS parameter 
    var port,conf string
    var err error
    flag.StringVar(&port, "exporter.port", "9104", "web Listen Port")
	flag.StringVar(&conf,"conf","config.yml","Configure YAML")
    flag.StringVar(&region,"region","ap-northeast-2","AWS Regions")
    flag.StringVar(&instance,"instance","","AWS Instance")
    flag.Parse()

    // Not Set Instance Exception
    if instance == "" {
        log.Errorf("Not set InstanceIdentifier.")
        os.Exit(1)
    }

    var b []byte
  
    if b, err = ioutil.ReadFile(conf); err != nil {
        log.Errorf("Failed to read config file: %s", err)
        os.Exit(1)
    }
  
    // Load yaml
    if err := yaml.Unmarshal(b, &config); err != nil {
        log.Errorf("Failed to load config: %s", err)
        os.Exit(1)
    }

    // Regist handler
	log.Infof("Regist version collector - %s", collector)
    prometheus.Register(version.NewCollector(collector))
	prometheus.Register(&CloudCollector{})
 
    // Regist http handler
    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        h := promhttp.HandlerFor(prometheus.Gatherers{
            prometheus.DefaultGatherer,
        }, promhttp.HandlerOpts{})
        h.ServeHTTP(w, r)
    })
 
    bind := fmt.Sprintf("0.0.0.0:%s",port)
    // start server
    log.Infof("Starting http server - %s", bind)
    if err := http.ListenAndServe(bind, nil); err != nil {
        log.Errorf("Failed to start http server: %s", err)
    }
}

type CloudCollector struct{}

// Describe Prometheus describe
func (e *CloudCollector) Describe(ch chan<- *prometheus.Desc){
    Labels := []string{"instance","statistics"}

    for metricName, metric := range config.Metrics {
        metric.metricDesc = prometheus.NewDesc(
            prometheus.BuildFQName("cloud_exporter", "", metricName),
            metric.Description,
            Labels, nil,
        )
        config.Metrics[metricName] = metric
        log.Infof("metric description for \"%s\" registerd", metricName)
    }
}

// Collect Prometheus collect
func (e *CloudCollector) Collect(ch chan<- prometheus.Metric){
    // Session
    cw := cloudwatch.New(session.New(), &aws.Config{
        Region: aws.String(region),
    })

    for name, metric := range config.Metrics {
        now 	:= time.Now()
        data 	:= now.UTC().Format(time.RFC3339)
        end, _ 	:= time.Parse(time.RFC3339,data)
        start 	:= end.Add(time.Minute * -5)

        vl, err := cw.GetMetricStatistics(
            &cloudwatch.GetMetricStatisticsInput{
                Namespace: 		aws.String("AWS/RDS"),
                MetricName: 	aws.String(name),
                Period: 		aws.Int64(config.Aws.Period),
                StartTime:		aws.Time(start),
                EndTime:		aws.Time(end),
                Statistics:     []*string{
                    aws.String(cloudwatch.StatisticAverage),
                    aws.String(cloudwatch.StatisticSum),
                    aws.String(cloudwatch.StatisticMinimum),
                    aws.String(cloudwatch.StatisticMaximum),
                },
                Dimensions: []*cloudwatch.Dimension{
                    {
                        Name:  aws.String("DBInstanceIdentifier"),
                        Value: aws.String(instance),
                    },
                },
            })
        if err != nil {
            panic(err)
        }

        // Last Time Data
        var lstIdx int = 0
        for i:=0; i<len(vl.Datapoints);i++ {
            if i != 0 {
                x := aws.TimeValue(vl.Datapoints[lstIdx].Timestamp)
                y := aws.TimeValue(vl.Datapoints[i].Timestamp)

                if x.After(y) {
                    lstIdx = lstIdx
                } else {
                    lstIdx = i
                }
            }
        }

         // Data Debug
         if config.Aws.Logmode == "debug" {
            fmt.Println(instance)
            fmt.Println(vl)
        }

        // Set Labels / values
        var avgLabels,minLabels,maxLabels,sumLabels []string 
        var avgVal,minVal,maxVal,sumVal float64

        // Average
        avgLabels = append(avgLabels,instance)
        avgLabels = append(avgLabels,"avg")
        // Minimum
        minLabels = append(minLabels,instance)
        minLabels = append(minLabels,"min")
        // Max
        maxLabels = append(maxLabels,instance)
        maxLabels = append(maxLabels,"max")
        // Sum
        sumLabels = append(sumLabels,instance)
        sumLabels = append(sumLabels,"sum")

         // Set Valeus
        if len(vl.Datapoints) == 0 {
            avgVal = 0
            minVal = 0
            maxVal = 0
            sumVal = 0
        } else {
            avgVal  = aws.Float64Value(vl.Datapoints[lstIdx].Average)
            minVal  = aws.Float64Value(vl.Datapoints[lstIdx].Minimum)
            maxVal  = aws.Float64Value(vl.Datapoints[lstIdx].Maximum)
            sumVal  = aws.Float64Value(vl.Datapoints[lstIdx].Sum)
        }
       
        // Channel Values
        ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, avgVal, avgLabels...)
        ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, minVal, minLabels...)
        ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, maxVal, maxLabels...)
        ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, sumVal, sumLabels...)

    }
}
   