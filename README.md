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
db.measures.createIndex( { "sensor.sensorId": 1,  "timestamp": -1 } );
```

```
db.measures.aggregate([
  { $sort: { "sensor.sensorId": 1, timestamp: -1 } },
  { $group: { _id: "$sensor.sensorId", timestamp: { $first: "$timestamp" }, value: { $first: "$value" } } }
]).toArray();
```