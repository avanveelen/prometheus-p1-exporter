package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/jordyv/prometheus-p1-exporter/conn"
	"github.com/jordyv/prometheus-p1-exporter/parser"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var readInterval time.Duration
var listenAddr string
var apiEndpoint string
var useMock bool
var verbose bool
var metricNamePrefix = "p1_"
var usbSerial string

var (
	registry                         = prometheus.NewRegistry()
	electricityConsumptionHighMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "consumption_electricity_high",
		Help: "1.8.1 - Electricity consumption high tariff in kWh",
	})
	electricityConsumptionLowMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "consumption_electricity_low",
		Help: "1.8.2 - Electricity consumption low tariff in kWh",
	})
	electricityProductionHighMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "production_electricity_high",
		Help: "2.8.1 - Electricity production high tariff in kWh",
	})
	electricityProductionLowMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "production_electricity_low",
		Help: "2.8.2 - Electricity production low tariff in kWh",
	})
	actualElectricityConsumptionMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "actual_electricity_consumption",
		Help: "1.7.0 - Actual electricity power consumption in kW",
	})
	actualElectricityProductionMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "actual_electricity_production",
		Help: "2.7.0 - Actual electricity power production in kW",
	})
	activeTariffMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "active_tariff",
		Help: "96.14.0 - Active tariff",
	})
	powerFailuresLongMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "power_failures_long",
		Help: "96.7.9 - Power failures long count",
	})
	powerFailuresShortMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: metricNamePrefix + "power_failures_short",
		Help: "96.7.21 - Power failures short count",
	})
	gasConsumptionMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Name: metricNamePrefix + "consumption_gas",
		Help: "24.2.1 - Gas usage in m³",
	})
)

func init() {
	registry.MustRegister(electricityConsumptionHighMetric)
	registry.MustRegister(electricityConsumptionLowMetric)
	registry.MustRegister(electricityProductionHighMetric)
	registry.MustRegister(electricityProductionLowMetric)
	registry.MustRegister(actualElectricityConsumptionMetric)
	registry.MustRegister(actualElectricityProductionMetric)
	registry.MustRegister(activeTariffMetric)
	registry.MustRegister(powerFailuresLongMetric)
	registry.MustRegister(powerFailuresShortMetric)
	registry.MustRegister(gasConsumptionMetric)
}

func main() {
	flag.StringVar(&listenAddr, "listen", "127.0.0.1:8888", "Listen address for HTTP metrics")
	flag.DurationVar(&readInterval, "interval", 10*time.Second, "Interval between metric reads")
	flag.BoolVar(&useMock, "mock", false, "Use dummy source instead of ttyUSB0 socket")
	flag.StringVar(&apiEndpoint, "apiEndpoint", "", "Use API endpoint to read the telegram (use for HomeWizard)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output logging")
	flag.StringVar(&usbSerial, "usbserial", "/dev/ttyUSB0", "USB serial device")
	flag.Parse()

	var source conn.Source
	if useMock {
		source = conn.NewMockSource()
	} else if apiEndpoint != "" {
		source = conn.NewAPISource(apiEndpoint)
	} else {
		source = conn.NewSerialSource(usbSerial)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	go func() {
		errorCount := 0

		// Initialize variables to store the last cumulative values
		var lastElectricityUsageHigh, lastElectricityUsageLow, lastElectricityReturnedHigh, lastElectricityReturnedLow, lastGasUsage float64

		for {
			if errorCount > 10 {
				logrus.Errorln("Quitting because there were too many errors")
				os.Exit(1)
			}

			lines, err := conn.ReadTelegram(&conn.ESMR5TelegramReaderOptions, source)
			if err != nil {
				logrus.Errorln("Error while reading telegram from source", err)
				errorCount++
				time.Sleep(readInterval)
				continue
			}
			telegram, err := parser.ParseTelegram(&parser.XS210ESMR5TelegramFormat, lines)
			if err != nil {
				logrus.Errorln("Error while parsing telegram", err)
				errorCount++
				time.Sleep(readInterval)
				continue
			}
			errorCount = 0

			// Update counters with differences from last readings for cumulative values
			if telegram.ElectricityUsageHigh != nil {
				diff := *telegram.ElectricityUsageHigh - lastElectricityUsageHigh
				if diff > 0 {
					electricityConsumptionHighMetric.Add(diff)
					lastElectricityUsageHigh = *telegram.ElectricityUsageHigh
				}
			}
			if telegram.ElectricityUsageLow != nil {
				diff := *telegram.ElectricityUsageLow - lastElectricityUsageLow
				if diff > 0 {
					electricityConsumptionLowMetric.Add(diff)
					lastElectricityUsageLow = *telegram.ElectricityUsageLow
				}
			}
			if telegram.ElectricityReturnedHigh != nil {
				diff := *telegram.ElectricityReturnedHigh - lastElectricityReturnedHigh
				if diff > 0 {
					electricityProductionHighMetric.Add(diff)
					lastElectricityReturnedHigh = *telegram.ElectricityReturnedHigh
				}
			}
			if telegram.ElectricityReturnedLow != nil {
				diff := *telegram.ElectricityReturnedLow - lastElectricityReturnedLow
				if diff > 0 {
					electricityProductionLowMetric.Add(diff)
					lastElectricityReturnedLow = *telegram.ElectricityReturnedLow
				}
			}
			if telegram.GasUsage != nil {
				diff := *telegram.GasUsage - lastGasUsage
				if diff > 0 {
					gasConsumptionMetric.Add(diff)
					lastGasUsage = *telegram.GasUsage
				}
			}

			// Update gauges directly for instantaneous values
			if telegram.ActualElectricityDelivered != nil {
				actualElectricityConsumptionMetric.Set(*telegram.ActualElectricityDelivered)
			}
			if telegram.ActualElectricityRetreived != nil {
				actualElectricityProductionMetric.Set(*telegram.ActualElectricityRetreived)
			}

			activeTariffMetric.Set(float64(telegram.ActiveTariff))
			powerFailuresLongMetric.Set(float64(telegram.PowerFailuresLong))
			powerFailuresShortMetric.Set(float64(telegram.PowerFailuresShort))

			logrus.Debugf("%+v\n", telegram)

			time.Sleep(readInterval)
		}
	}()

	logrus.Infoln("Start listening at", listenAddr)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	logrus.Fatalln(http.ListenAndServe(listenAddr, nil))
}
