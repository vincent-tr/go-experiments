# go-experiments

## run

### collector

```bash
cd mylife-energy/
go run cmd/collector/*.go
```

### web

```bash
cd mylife-energy/
go run cmd/web/*.go
```

## docker build

### collector

```bash
docker build -t go-mylife-energy-collector -f docker/collector/Dockerfile .
```

### web

```bash
docker build -t go-mylife-energy-web -f docker/web/Dockerfile .
```

TODO: one image with command switch for collector/web (like gallery?)

```
db.measures.aggregate([
                     { $match: { "sensor.sensorId": {$regex : "real"} } },
                     { $sort: { timestamp: 1 } },
                     { $group: { _id: "$sensor.sensorId", timestamp: { $last: "$timestamp" }, value: { $last: "$value" } } },
                     { $sort: { _id: 1 } }
                   ]).toArray()
```