package store

import "mylife-tools-server/services/database"

type CollectionBuilder[TEntity Entity] struct {
	colName     string
	dataColName string
	decoder     func(raw []byte) (TEntity, error)
	encoder     func(obj TEntity) ([]byte, error)
}

type icollectionBuilder interface {
	name() string
	build() genericInternalCollection
}

func MakeCollectionBuilder[TEntity Entity](
	colName string,
	dataColName string,
	decoder func(raw []byte) (TEntity, error),
	encoder func(obj TEntity) ([]byte, error),
) *CollectionBuilder[TEntity] {
	return &CollectionBuilder[TEntity]{
		colName:     colName,
		dataColName: dataColName,
		decoder:     decoder,
		encoder:     encoder,
	}
}

func (builder *CollectionBuilder[TEntity]) name() string {
	return builder.colName
}

func (builder *CollectionBuilder[TEntity]) build() genericInternalCollection {
	col := &collection[TEntity]{
		container: NewContainer[TEntity](builder.colName),
	}

	dataCol := database.GetCollection(builder.dataColName)
	col.updater = makeUpdater[TEntity](col, dataCol, builder.decoder, builder.encoder)

	return col
}
