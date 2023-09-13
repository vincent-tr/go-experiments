# go-experiments

type State[T] interface // mylife:`name=toto,description=tata,type=titi`

Get
Set

type Action[T] interface // mylife:`name=toto,description=tata,type=titi`

RegisterCallback

type Config[T] interface // mylife:`type=titi`,

Get

## generate

```shell
go generate mylife-home-core-plugins-logic-base/main.go 
```

## run

```shell
go run mylife-home-core/main.go server
```