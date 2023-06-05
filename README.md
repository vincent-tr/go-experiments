# go-experiments

```
cd mylife-energy/
go run cmd/main.go
```

```
db.measures.aggregate([
                     { $match: { "sensor.sensorId": {$regex : "real"} } },
                     { $sort: { timestamp: 1 } },
                     { $group: { _id: "$sensor.sensorId", timestamp: { $last: "$timestamp" }, value: { $last: "$value" } } },
                     { $sort: { _id: 1 } }
                   ]).toArray()
```