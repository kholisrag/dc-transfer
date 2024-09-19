# MongoDB Source Connector Documentation

## Overview

The MongoDB Source Connector enables ingestion of data from MongoDB into the target system. It supports a variety of ingestion modes, including snapshot-based ingestion and real-time data capture using the MongoDB oplog. The connector can be customized with a range of options for managing collections, authentication, replication, and data filtering.

It supports various ingestion modes, including:

- **Snapshot**: Fetches all data from a specified table.
- **Snapshot with Cursor Column**: Fetches data in recurring intervals with a filter.
- **Realtime CDC (Change Data Capture)**: Captures ongoing changes to the database in real time.

Source connector ingest semi-structured data. All collections ingested as `Document` table structure:


This document provides a detailed explanation of the configuration options for the MongoDB Source Connector, which are controlled via JSON or YAML formats using the `MongoSource` Go structure.


## Data Structure

Data ingested from MongoDB is structured as a table with two columns:

- **ID**: A string field representing the MongoDB document's `_id` field. The `_id` is a unique identifier for each MongoDB document.
    - **Type**: String
    - **Primary Key**: True
    - **Original Type**: `mongo:bson_id`

- **Document**: A field containing the full BSON document from MongoDB.
    - **Type**: Any (Serialized JSON or BSON)
    - **Original Type**: `mongo:bson`

This basic structure can be flattened or transformed as needed by the target system. The engine supports further transformation for future data flattening. This flexibility allows you to extract or restructure nested fields within the `Document` for more efficient storage and querying in downstream systems.

This two-column schema ensures that each MongoDB document is ingested in its entirety, with the option to flatten the document for efficient querying in downstream systems.



---

## Configuration

The MongoDB Source Connector can be configured using the `MongoSource` structure. Below is a breakdown of each configuration field.

### JSON/YAML Example

```json
{
  "Hosts": ["mongo-host1", "mongo-host2"],
  "Port": 27017,
  "ReplicaSet": "replica-set-name",
  "AuthSource": "admin",
  "User": "mongo-user",
  "Password": "mongo-password",
  "Collections": [
    {
      "Database": "db1",
      "Collection": "coll1"
    }
  ],
  "ExcludedCollections": [
    {
      "Database": "db2",
      "Collection": "temp_coll"
    }
  ],
  "TechnicalDatabase": "",
  "IsHomo": false,
  "SlotID": "transfer-id",
  "SecondaryPreferredMode": true,
  "TLSFile": "FILE CONTENT",
  "ReplicationSource": "Oplog",
  "BatchingParams": {
    "MaxBatchSize": 10000,
    "FlushInterval": "5s"
  },
  "DesiredPartSize": 50000,
  "PreventJSONRepack": false,
  "FilterOplogWithRegexp": true,
  "Direct": false,
  "RootCAFiles": ["/path/to/ca1.pem", "/path/to/ca2.pem"],
  "SRVMode": true
}
```

### Fields

- **Hosts** (`[]string`): List of MongoDB hosts in the replica set or sharded cluster. For example: `["host1", "host2"]`.

- **Port** (`int`): The port number for MongoDB connections. Default is typically `27017`.

- **ReplicaSet** (`string`): The name of the MongoDB replica set. This ensures the connector connects to a MongoDB replica set for better redundancy and performance.

- **AuthSource** (`string`): Specifies the authentication database (commonly `admin`). This database is used for authentication purposes.

- **User** (`string`): The MongoDB username for connecting to the database.

- **Password** (`server.SecretString`): The password for the MongoDB user.

- **Collections** (`[]MongoCollection`): Specifies the collections to be included in the data ingestion. Each entry is a combination of the `Database` and `Collection` fields.

    - Example:
      ```json
      {
        "Database": "db1",
        "Collection": "coll1"
      }
      ```

- **ExcludedCollections** (`[]MongoCollection`): Specifies the collections to be excluded from the ingestion process.

- **TechnicalDatabase** (`string`, deprecated): This field is deprecated and should always be left empty (`""`).

- **IsHomo** (`bool`): Specifies whether homogeneous transfers are enabled, typically used in homogeneous environments like MongoDB to MongoDB replication.

- **SlotID** (`string`): A synthetic entity used for data replication. This is always equal to the `transfer_id` and helps manage replication slots.

- **SecondaryPreferredMode** (`bool`): If `true`, the connector prefers to read from MongoDB secondary nodes in the replica set for better load distribution.

- **TLSFile** (`string`): Path to the TLS certificate for secure connections to MongoDB.

- **ReplicationSource** (`MongoReplicationSource`): Specifies the replication mode. The two common modes are:
    - **Oplog**: Real-time change data capture using the MongoDB oplog.
    - **Snapshot**: Captures a snapshot of the data at a specific point in time.

- **BatchingParams** (`BatcherParameters`): Internal parameters for batching during data ingestion. These control the size of batches and flush intervals:
    - **MaxBatchSize**: Maximum number of documents per batch.
    - **FlushInterval**: Time interval for flushing data in the batch.

- **DesiredPartSize** (`uint64`): The desired size of each data part (in terms of documents) when ingesting data. This helps with parallelism and performance optimization.

- **PreventJSONRepack** (`bool`): Controls whether JSON repacking is prevented during ingestion. This field is typically used during migrations and specific use cases where performance might be impacted by JSON manipulation.

- **FilterOplogWithRegexp** (`bool`): When enabled (`true`), the MongoDB oplog is filtered using regular expressions that match the specified collections and excluded collections. This can improve network efficiency but may result in lost oplog events if no changes occur in the listened-to collections.

    - **Advantages**: Better network efficiency.
    - **Disadvantages**: Potential for data loss if no changes occur in the monitored collections.

- **Direct** (`bool`): When set to `true`, the connector makes a direct connection to the MongoDB instance, bypassing the connection string handling for `mongodb+srv`. Refer to the [MongoDB Go Driver connection guide](https://www.mongodb.com/docs/drivers/go/current/fundamentals/connections/connection-guide/) for more details.

- **RootCAFiles** (`[]string`): List of paths to root CA files for verifying SSL connections to MongoDB.

- **SRVMode** (`bool`): When set to `true`, the MongoDB client uses the `mongodb+srv` connection format. This is commonly used for cloud-based MongoDB services.

---

## Ingestion Modes

### 1. Snapshot Ingestion

In snapshot mode, the connector captures a static view of the MongoDB collections at a specific point in time.

- **Use Case**: Initial data loads or one-time migrations.
- **Configuration**: No special settings are required for snapshot mode. The connector will automatically handle the full data copy from the specified collections.

### 2. Real-Time Change Data Capture (Oplog Replication)

In oplog mode, the connector listens to changes in MongoDB using the oplog. This enables real-time ingestion of data updates as they happen in the source MongoDB database.

- **Use Case**: Real-time synchronization of data between MongoDB and the target system.
- **Configuration**: Set `ReplicationSource` to `Oplog` and configure the connector to read from the oplog.

#### Oplog Filtering with Regular Expressions

You can filter the oplog using regular expressions by setting `FilterOplogWithRegexp` to `true`. This will limit the oplog events to the specified collections and exclude unwanted collections from being ingested.

- **Use Case**: When you want to ingest data only from specific collections to reduce network and processing load.

---

## Secure Connections (TLS)

To secure the connection to MongoDB, you can configure the `TLSFile` field to provide the path to the TLS certificate. Additionally, you can use the `RootCAFiles` field to specify the paths to root CA files for validating SSL certificates.

```yaml
TLSFile: "FILE_CONTENT"
RootCAFiles:
  - "/path/to/ca1.pem"
  - "/path/to/ca2.pem"
```

---

## Advanced Features

### SecondaryPreferredMode

By setting `SecondaryPreferredMode` to `true`, the connector can read data from secondary nodes in the MongoDB replica set. This reduces the load on the primary node and improves overall performance.

```yaml
SecondaryPreferredMode: true
```

### Direct Connections

For environments where direct connections are preferred over `mongodb+srv` URIs, you can enable the `Direct` field to force a direct connection to the MongoDB host.

```yaml
Direct: true
```

### Collection Inclusion and Exclusion

You can specify the collections to include or exclude using the `Collections` and `ExcludedCollections` fields, respectively. This allows you to control precisely which data is ingested from MongoDB.

- **Collections**: Only these collections will be ingested.
- **ExcludedCollections**: These collections will be skipped during ingestion.

```yaml
Collections:
  - Database: "db1"
    Collection: "coll1"
ExcludedCollections:
  - Database: "db2"
    Collection: "temp_coll"
```

---

## Error Handling and Deprecated Features

1. **TechnicalDatabase**: This field is deprecated and should always be set to an empty string (`""`).

2. **Replication Consistency**: Ensure that the `SlotID` is properly configured (it should always match the `transfer_id`), especially in replication scenarios. This prevents data inconsistency during the replication process.

3. **PreventJSONRepack**: This option is generally not recommended for use except in very specific migration scenarios where JSON repacking is known to cause performance issues.

---

## Demo

TODO

