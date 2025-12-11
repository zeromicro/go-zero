# Migrating from 1.x to 2.0

To upgrade imports of the Go Driver from v1 to v2, we recommend using [marwan-at-work/mod
](https://github.com/marwan-at-work/mod):

```
mod upgrade --mod-name=go.mongodb.org/mongo-driver
```

# Notice
After completing the mod upgrade, code changes are typically unnecessary in the vast majority of cases. However, if your project references packages including but not limited to those listed below, you'll need to manually replace them, as these libraries are no longer present in the v2 version.
```go
go.mongodb.org/mongo-driver/bson/bsonrw    => go.mongodb.org/mongo-driver/v2/bson
go.mongodb.org/mongo-driver/bson/bsoncodec => go.mongodb.org/mongo-driver/v2/bson
go.mongodb.org/mongo-driver/bson/primitive => go.mongodb.org/mongo-driver/v2/bson
```   

See the following resources to learn more about upgrading from version 1.x to 2.0.:
https://raw.githubusercontent.com/mongodb/mongo-go-driver/refs/heads/master/docs/migration-2.0.md