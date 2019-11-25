package main

import (
	"bufio"
	"net/http"
	"strconv"

	"github.com/roaldnefs/go-dsmr"

	"github.com/tarm/serial"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gopkg.in/alecthomas/kingpin.v2"

	log "github.com/sirupsen/logrus"
)

var (
	powerDelivered = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_power_delivered_kw",
		Help: "Actual electricity power delivered (+P) in 1 Watt resolution.",
	})
	instVoltL1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_voltage_l1_v",
		Help: "Instantaneous voltage L1 in V resolution.",
	})
	instVoltL2 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_voltage_l2_v",
		Help: "Instantaneous voltage L2 in V resolution.",
	})
	instVoltL3 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_voltage_l3_v",
		Help: "Instantaneous voltage L3 in V resolution.",
	})
	instCurL1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_current_l1_a",
		Help: "Instantaneous current L1 in A resolution.",
	})
	instCurL2 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_current_l2_a",
		Help: "Instantaneous current L1 in A resolution.",
	})
	instCurL3 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "p1_electricity_instantaneous_current_l3_a",
		Help: "Instantaneous current L1 in A resolution.",
	})
	config *serial.Config
)

func recordMetrics() {
	//	var channel = make(chan []byte)
	powerDelivered.Set(float64(0))
	instVoltL1.Set(float64(0))

	go func() {
		// Open serial port
		stream, err := serial.OpenPort(config)
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(stream)

		for {
			// Peek at the next byte, and look for the start of the telegram
			if peek, err := reader.Peek(1); err == nil {
				// The telegram starts with a '/' character keep reading
				// bytes until the start of the telegram is found
				if string(peek) != "/" {
					reader.ReadByte()
					continue
				}
			} else {
				continue
			}

			// Keep reading until the '!' character which indicates the end of
			// the telegram and is followed by the CRC
			rawTelegram, err := reader.ReadBytes('!')
			if err != nil {
				log.Error(err)
				continue
			}

			// Read the CRC which can be used to detect faulty telegram
			// TODO check CRC
			_, err = reader.ReadBytes('\n')
			if err != nil {
				log.Error(err)
				continue
			}

			telegram, err := dsmr.ParseTelegram(string(rawTelegram))
			if err != nil {
				log.Error(err)
				continue
			}

			if rawValue, ok := telegram.InstantaneousVoltageL1(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instVoltL1.Set(value)
			}

			if rawValue, ok := telegram.InstantaneousVoltageL2(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instVoltL2.Set(value)
			}

			if rawValue, ok := telegram.InstantaneousVoltageL3(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instVoltL3.Set(value)
			}

			if rawValue, ok := telegram.InstantaneousCurrentL1(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instCurL1.Set(value)
			}

			if rawValue, ok := telegram.InstantaneousCurrentL2(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instCurL2.Set(value)
			}

			if rawValue, ok := telegram.InstantaneousCurrentL3(); ok {
				value, err := strconv.ParseFloat(rawValue, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				instCurL3.Set(value)
			}

			i, ok := telegram.DataObjects["1-0:1.7.0"]
			if ok {
				value, err := strconv.ParseFloat(i.Value, 64)
				if err != nil {
					log.Error(err)
					continue
				}
				powerDelivered.Set(value)
			}
		}
	}()

}

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9602").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		serialPort = kingpin.Flag(
			"serial.port",
			"Serial port for the connection to the P1 interface.",
		).Required().String()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Serial configuration
	config = &serial.Config{
		Name: *serialPort,
		Baud: 115200,
	}

	registry := prometheus.NewRegistry()

	registry.MustRegister(powerDelivered)
	registry.MustRegister(instVoltL1)
	registry.MustRegister(instVoltL2)
	registry.MustRegister(instVoltL3)
	registry.MustRegister(instCurL1)
	registry.MustRegister(instCurL2)
	registry.MustRegister(instCurL3)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	log.WithFields(log.Fields{
		"version": "unknown",
	}).Info("Starting P1 Exporter")

	recordMetrics()

	http.Handle(*metricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
							<head><title>P1 Exporter</title></head>
							<body>
							<h1>P1 Exporter</h1>
							<p><a href="` + *metricsPath + `">Metrics</a></p>
							</body>
							</html>`))
	})

	log.WithFields(log.Fields{
		"listen_address": *listenAddress,
		"metrics_path":   *metricsPath,
	}).Info("Listing on " + *listenAddress)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
