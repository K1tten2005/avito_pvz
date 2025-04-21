package metrics

import (
  "github.com/prometheus/client_golang/prometheus"
)

type ProductMetrics struct {
  ReceptionTotal      prometheus.Counter
  PvzTotal     prometheus.Counter
  ProductTotal prometheus.Counter
}

func NewProductMetrics() (*ProductMetrics, error) {
  var metr ProductMetrics
  metr.ReceptionTotal = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "ReceptionTotal",
      Help: "Number of total hits.",
    })
  if err := prometheus.Register(metr.ReceptionTotal); err != nil {
    return nil, err
  }
  metr.PvzTotal = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "PvzTotal",
      Help: "Number of total hits.",
    })
  if err := prometheus.Register(metr.PvzTotal); err != nil {
    return nil, err
  }
  metr.ProductTotal = prometheus.NewCounter(
    prometheus.CounterOpts{
      Name: "ProductTotal",
      Help: "Number of total hits.",
    })
  if err := prometheus.Register(metr.ProductTotal); err != nil {
    return nil, err
  }

  return &metr, nil
}
func (m *ProductMetrics) IncreaseReceptionTotal() {
  m.ReceptionTotal.Inc()
}
func (m *ProductMetrics) IncreasePvzTotal() {
	m.PvzTotal.Inc()
}
func (m *ProductMetrics) IncreaseProductTotal() {
	m.PvzTotal.Inc()
}