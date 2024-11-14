# Prometheus P1 exporter #

Prometheus exporter for smart meter statistics fetched with a P1 cable.   
This fork has improved prometheus metric names.


## Installation ##

### From source ###

With Go get:

```
$ go get github.com/jordyv/prometheus-p1-exporter
```

With Go install (recommended):

```
$ go install github.com/jordyv/prometheus-p1-exporter@latest

```

Make:

```
$ git clone https://github.com/jordyv/prometheus-p1-exporter.git
$ cd prometheus-p1-exporter
$ make
$ sudo make install
```

## Usage ##

```
Usage of ./prometheus-p1-exporter:
  -apiEndpoint string
        Use API endpoint to read the telegram (use for HomeWizard)
  -interval duration
        Interval between metric reads (default 10s)
  -listen string
        Listen address for HTTP metrics (default "127.0.0.1:8888")
  -mock
        Use dummy source instead of ttyUSB0 socket
  -usbserial string
    	USB serial device (default "/dev/ttyUSB0")
  -verbose
        Verbose output logging
```

By default the exporter will collect metrics from `/dev/ttyUSB0` every 10 seconds and export the metrics to an HTTP endpoint at `http://127.0.0.1:8888/metrics`. This endpoint can be added to your Prometheus configuration.

Example metrics page:

```
# HELP p1_active_tariff 96.14.0 - Active tariff
# TYPE p1_active_tariff gauge
p1_active_tariff 1
# HELP p1_actual_electricity_consumption 1.7.0 - Actual electricity power consumption in kW
# TYPE p1_actual_electricity_consumption gauge
p1_actual_electricity_consumption 0.279
# HELP p1_actual_electricity_production 2.7.0 - Actual electricity power production in kW
# TYPE p1_actual_electricity_production gauge
p1_actual_electricity_production 0
# HELP p1_consumption_electricity_high 1.8.1 - Electricity consumption high tariff in kWh
# TYPE p1_consumption_electricity_high counter
p1_consumption_electricity_high 10878.601
# HELP p1_consumption_electricity_low 1.8.2 - Electricity consumption low tariff in kWh
# TYPE p1_consumption_electricity_low counter
p1_consumption_electricity_low 10026.591
# HELP p1_consumption_gas 24.2.1 - Gas usage in m³
# TYPE p1_consumption_gas counter
p1_consumption_gas 5262.398
# HELP p1_power_failures_long 96.7.9 - Power failures long count
# TYPE p1_power_failures_long gauge
p1_power_failures_long 3
# HELP p1_power_failures_short 96.7.21 - Power failures short count
# TYPE p1_power_failures_short gauge
p1_power_failures_short 1
# HELP p1_production_electricity_high 2.8.1 - Electricity production high tariff in kWh
# TYPE p1_production_electricity_high counter
p1_production_electricity_high 359.608
# HELP p1_production_electricity_low 2.8.2 - Electricity production low tariff in kWh
# TYPE p1_production_electricity_low counter
p1_production_electricity_low 152.981
```

## Development ##

Currently only the DSMR 5.0 format is supported and the parser is default configured to parse the telegram message with the keys the Sagemcom XS210 is using.
If you have to support a different DSMR 5.0 message, feel free to create your own implementation of the TelegramFormat struct. To support a different format then DSMR 5.0 you can implement your own implementation of the TelegramReaderOptions struct.
