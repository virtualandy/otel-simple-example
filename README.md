## What is it?

The otel dice demo from https://opentelemetry.io/docs/languages/go/getting-started/

## Getting Started

Follow the instructions at [the OpenTelemetry go page](https://opentelemetry.io/docs/languages/go/getting-started/#setup).

## Using this Repo

Run the `TBD` branch to run without any otel.
Run the `otel-with-stdout` branch to run and have the otel exported to stdout

### Visualizing Traces with Jaeger

Run the `otel-with-viz` branch to see traces in Jaeger.

First, run the docker container for Jaeger:

```
docker run -d --name jaeger \
-e COLLECTOR_OTLP_ENABLED=true \
-p 16686:16686 \
-p 4317:4317 \
-p 4318:4318 \ jaegertracing/all-in-one:latest
```

Then, build + run the code:

```
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 go run main.go rolldice.go otel.go
```

_Note:_ You'll need `OTEL_EXPORTER_OTLP_ENDPOINT` pointed to the Jaeger container. The [default for golang](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp) is `https://localhost:4318` which won't work if you run jaeger as above since it's listening on plain ol' `http://`.

Todo:

- [ ] Create a branch same code but with basic logs instead of OTel
- [ ] Create a branch (main) with OTel
- [x] Add Jaeger or Honeycomb or a visualizer
