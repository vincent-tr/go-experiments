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

## Setup/Notes

### Mongo

index for live query
```
db.measures.createIndex( { "sensor.sensorId": 1,  "timestamp": -1 } );
```

query
```
db.measures.aggregate([
  { $sort: { "sensor.sensorId": 1, timestamp: -1 } },
  { $group: { _id: "$sensor.sensorId", timestamp: { $first: "$timestamp" }, value: { $first: "$value" } } }
]).toArray();
```

### Tesla API

get token : (work only using master branch currently)
```
git clone https://github.com/bogosj/tesla.git
cd tesla
go get
cd cmd/login
go work use ../..
go run . -o ~/tesla.token
```