package query

import "go.mongodb.org/mongo-driver/bson"

type PipelineBuilder struct {
	stages []bson.D
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		stages: make([]bson.D, 0),
	}
}

func (builder *PipelineBuilder) Build() interface{} {
	return builder.stages
}

func (builder *PipelineBuilder) addStage(name string, value bson.D) *PipelineBuilder {
	stage := bson.D{{Key: name, Value: value}}
	builder.stages = append(builder.stages, stage)

	return builder
}

type SortOrder int

const (
	Asc  SortOrder = 1
	Desc SortOrder = -1
)

type SortField struct {
	Name  string
	Order SortOrder
}

func (builder *PipelineBuilder) Sort(fields ...SortField) *PipelineBuilder {
	value := bson.D{}

	for _, item := range fields {
		value = append(value, bson.E{Key: item.Name, Value: int(item.Order)})
	}

	return builder.addStage("$sort", value)
}

type GroupField struct {
	Name        string
	Accumulator string
	Expression  string
}

func (builder *PipelineBuilder) Group(groupKeyExpression string, fields ...GroupField) *PipelineBuilder {
	value := bson.D{
		{Key: "_id", Value: groupKeyExpression},
	}

	for _, item := range fields {
		value = append(value, bson.E{Key: item.Name, Value: bson.D{{Key: item.Accumulator, Value: item.Expression}}})
	}

	return builder.addStage("$group", value)
}
