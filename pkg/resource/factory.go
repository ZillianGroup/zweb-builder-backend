// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"github.com/illa-family/builder-backend/pkg/plugins/common"
	"github.com/illa-family/builder-backend/pkg/plugins/elasticsearch"
	"github.com/illa-family/builder-backend/pkg/plugins/mongodb"
	"github.com/illa-family/builder-backend/pkg/plugins/mysql"
	"github.com/illa-family/builder-backend/pkg/plugins/postgresql"
	"github.com/illa-family/builder-backend/pkg/plugins/redis"
	"github.com/illa-family/builder-backend/pkg/plugins/restapi"
)

var (
	REST_RESOURCE          = "restapi"
	MYSQL_RESOURCE         = "mysql"
	MARIADB_RESOURCE       = "mariadb"
	TIDB_RESOURCE          = "tidb"
	POSTGRES_RESOURCE      = "postgresql"
	REDIS_RESOURCE         = "redis"
	MONGODB_RESOURCE       = "mongodb"
	ELASTICSEARCH_RESOURCE = "elasticsearch"
)

type AbstractResourceFactory interface {
	Build() common.DataConnector
}

type Factory struct {
	Type string
}

func (f *Factory) Generate() common.DataConnector {
	switch f.Type {
	case REST_RESOURCE:
		restapiRsc := &restapi.RESTAPIConnector{}
		return restapiRsc
	case MYSQL_RESOURCE, MARIADB_RESOURCE, TIDB_RESOURCE:
		sqlRsc := &mysql.MySQLConnector{}
		return sqlRsc
	case POSTGRES_RESOURCE:
		pgsRsc := &postgresql.Connector{}
		return pgsRsc
	case REDIS_RESOURCE:
		redisRsc := &redis.Connector{}
		return redisRsc
	case MONGODB_RESOURCE:
		mongoRsc := &mongodb.Connector{}
		return mongoRsc
	case ELASTICSEARCH_RESOURCE:
		esRsc := &elasticsearch.Connector{}
		return esRsc
	default:
		return nil
	}
}
