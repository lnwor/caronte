/*
 * This file is part of caronte (https://github.com/eciavatta/caronte).
 * Copyright (c) 2020 Emiliano Ciavatta.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type TestStorageWrapper struct {
	DbName     string
	Storage    *MongoStorage
	Context    context.Context
	CancelFunc context.CancelFunc
}

func NewTestStorageWrapper(t *testing.T) *TestStorageWrapper {
	mongoHost, ok := os.LookupEnv("MONGO_HOST")
	if !ok {
		mongoHost = "localhost"
	}
	mongoPort, ok := os.LookupEnv("MONGO_PORT")
	if !ok {
		mongoPort = "27017"
	}
	port, err := strconv.Atoi(mongoPort)
	require.NoError(t, err, "invalid port")

	dbName := "caronte_test"

	storage, err := NewMongoStorage(mongoHost, port, dbName, "", "")
	require.NoError(t, err)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)

	require.NoError(t, storage.database.Drop(ctx))

	return &TestStorageWrapper{
		DbName:     dbName,
		Storage:    storage,
		Context:    ctx,
		CancelFunc: cancelFunc,
	}
}

func (tsw TestStorageWrapper) AddCollection(collectionName string) {
	tsw.Storage.collections[collectionName] = tsw.Storage.client.Database(tsw.DbName).Collection(collectionName)
}

func (tsw TestStorageWrapper) Destroy(t *testing.T) {
	err := tsw.Storage.client.Disconnect(tsw.Context)
	tsw.CancelFunc()
	require.NoError(t, err, "failed to disconnect to database")
}
