package resourcelist

var (
	TYPE_TRANSFORMER   = "transformer"
	TYPE_RESTAPI       = "restapi"
	TYPE_GRAPHQL       = "graphql"
	TYPE_REDIS         = "redis"
	TYPE_MYSQL         = "mysql"
	TYPE_MARIADB       = "mariadb"
	TYPE_POSTGRESQL    = "postgresql"
	TYPE_MONGODB       = "mongodb"
	TYPE_TIDB          = "tidb"
	TYPE_ELASTICSEARCH = "elasticsearch"
	TYPE_S3            = "s3"
	TYPE_SMTP          = "smtp"
	TYPE_SUPABASEDB    = "supabasedb"
	TYPE_FIREBASE      = "firebase"
	TYPE_CLICKHOUSE    = "clickhouse"
	TYPE_MSSQL         = "mssql"
	TYPE_HUGGINGFACE   = "huggingface"
	TYPE_DYNAMODB      = "dynamodb"
	TYPE_SNOWFLAKE     = "snowflake"
	TYPE_COUCHDB       = "couchdb"
	TYPE_HFENDPOINT    = "hfendpoint"
	TYPE_ORACLE        = "oracle"
	TYPE_APPWRITE      = "appwrite"
	TYPE_GOOGLESHEETS  = "googlesheets"
	TYPE_NEON          = "neon"
	TYPE_UPSTASH       = "upstash"
	TYPE_AIRTABLE      = "airtable"
	TYPE_HYDRA         = "hydra"
)

var (
	TYPE_TRANSFORMER_ID   = 0
	TYPE_RESTAPI_ID       = 1
	TYPE_GRAPHQL_ID       = 2
	TYPE_REDIS_ID         = 3
	TYPE_MYSQL_ID         = 4
	TYPE_MARIADB_ID       = 5
	TYPE_POSTGRESQL_ID    = 6
	TYPE_MONGODB_ID       = 7
	TYPE_TIDB_ID          = 8
	TYPE_ELASTICSEARCH_ID = 9
	TYPE_S3_ID            = 10
	TYPE_SMTP_ID          = 11
	TYPE_SUPABASEDB_ID    = 12
	TYPE_FIREBASE_ID      = 13
	TYPE_CLICKHOUSE_ID    = 14
	TYPE_MSSQL_ID         = 15
	TYPE_HUGGINGFACE_ID   = 16
	TYPE_DYNAMODB_ID      = 17
	TYPE_SNOWFLAKE_ID     = 18
	TYPE_COUCHDB_ID       = 19
	TYPE_HFENDPOINT_ID    = 20
	TYPE_ORACLE_ID        = 21
	TYPE_APPWRITE_ID      = 22
	TYPE_GOOGLESHEETS_ID  = 23
	TYPE_NEON_ID          = 24
	TYPE_UPSTASH_ID       = 25
	TYPE_AIRTABLE_ID      = 26
	TYPE_HYDRA_ID         = 27
)

var type_array = [28]string{
	TYPE_TRANSFORMER,
	TYPE_RESTAPI,
	TYPE_GRAPHQL,
	TYPE_REDIS,
	TYPE_MYSQL,
	TYPE_MARIADB,
	TYPE_POSTGRESQL,
	TYPE_MONGODB,
	TYPE_TIDB,
	TYPE_ELASTICSEARCH,
	TYPE_S3,
	TYPE_SMTP,
	TYPE_SUPABASEDB,
	TYPE_FIREBASE,
	TYPE_CLICKHOUSE,
	TYPE_MSSQL,
	TYPE_HUGGINGFACE,
	TYPE_DYNAMODB,
	TYPE_SNOWFLAKE,
	TYPE_COUCHDB,
	TYPE_HFENDPOINT,
	TYPE_ORACLE,
	TYPE_APPWRITE,
	TYPE_GOOGLESHEETS,
	TYPE_NEON,
	TYPE_UPSTASH,
	TYPE_AIRTABLE,
	TYPE_HYDRA,
}

var type_map = map[string]int{
	TYPE_TRANSFORMER:   TYPE_TRANSFORMER_ID,
	TYPE_RESTAPI:       TYPE_RESTAPI_ID,
	TYPE_GRAPHQL:       TYPE_GRAPHQL_ID,
	TYPE_REDIS:         TYPE_REDIS_ID,
	TYPE_MYSQL:         TYPE_MYSQL_ID,
	TYPE_MARIADB:       TYPE_MARIADB_ID,
	TYPE_POSTGRESQL:    TYPE_POSTGRESQL_ID,
	TYPE_MONGODB:       TYPE_MONGODB_ID,
	TYPE_TIDB:          TYPE_TIDB_ID,
	TYPE_ELASTICSEARCH: TYPE_ELASTICSEARCH_ID,
	TYPE_S3:            TYPE_S3_ID,
	TYPE_SMTP:          TYPE_SMTP_ID,
	TYPE_SUPABASEDB:    TYPE_SUPABASEDB_ID,
	TYPE_FIREBASE:      TYPE_FIREBASE_ID,
	TYPE_CLICKHOUSE:    TYPE_CLICKHOUSE_ID,
	TYPE_MSSQL:         TYPE_MSSQL_ID,
	TYPE_HUGGINGFACE:   TYPE_HUGGINGFACE_ID,
	TYPE_DYNAMODB:      TYPE_DYNAMODB_ID,
	TYPE_SNOWFLAKE:     TYPE_SNOWFLAKE_ID,
	TYPE_COUCHDB:       TYPE_COUCHDB_ID,
	TYPE_HFENDPOINT:    TYPE_HFENDPOINT_ID,
	TYPE_ORACLE:        TYPE_ORACLE_ID,
	TYPE_APPWRITE:      TYPE_APPWRITE_ID,
	TYPE_GOOGLESHEETS:  TYPE_GOOGLESHEETS_ID,
	TYPE_NEON:          TYPE_NEON_ID,
	TYPE_UPSTASH:       TYPE_UPSTASH_ID,
	TYPE_AIRTABLE:      TYPE_AIRTABLE_ID,
	TYPE_HYDRA:         TYPE_HYDRA_ID,
}

func GetResourceIDMappedType(id int) string {
	return type_array[id]
}

func GetResourceNameMappedID(name string) int {
	return type_map[name]
}