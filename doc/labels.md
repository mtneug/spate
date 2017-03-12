# List of `spate` labels

All `spate` labels are prefixed by `de.mtneug.spate`, which is therefore omitted in the following table.

| Label                                 | Description                                                |
| ------------------------------------- | ---------------------------------------------------------- |
| `autoscaler.period`                   | Autoscaling decision period.                               |
| `autoscaler.cooldown.scaled_up`       | Cooldown time after a scale up.                            |
| `autoscaler.cooldown.scaled_down`     | Cooldown time after a scale up.                            |
| `autoscaler.cooldown.service_added`   | Cooldown time after a service is added to `spate`.         |
| `autoscaler.cooldown.service_updated` | Cooldown time after a service is updated.                  |
| `replica.min`                         | Minimum number of replicas.                                |
| `replica.max`                         | Maximum number of replicas.                                |
| `metric.<key>.type`                   | Type of metric `<key>`.                                    |
| `metric.<key>.kind`                   | Kind of metric `<key>`.                                    |
| `metric.<key>.prometheus.endpoint`    | Prometheus endpoint URL of metric `<key>`.                 |
| `metric.<key>.prometheus.name`        | Prometheus metric name of metric `<key>`.                  |
| `metric.<key>.aggregation.method`     | Aggregation method for metric `<key>`.                     |
| `metric.<key>.aggregation.amount`     | Number of measurements to aggregate for metric `<key>`.    |
| `metric.<key>.observer.period`        | Measurement period of metric `<key>`.                      |
| `metric.<key>.target`                 | Target value of metric `<key>`.                            |
| `metric.<key>.target.deviation.lower` | Allowed lower deviation from the target of metric `<key>`. |
| `metric.<key>.target.deviation.upper` | Allowed upper deviation from the target of metric `<key>`. |
