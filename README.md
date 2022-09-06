# cloud_exporter
## Description
cloud_exporter는 DBA 성동찬님이 개발하신 Query Exporter 코드를 응용하여 개발된 Exporter 입니다. AWS CloudWatch Log를 수집하고, 원하는 메트릭만 선택하여 수집할 수 있도록 개발되었습니다. 

## Configure
---
```
AWS:
  period: 60
  logmode: debug ## info/debug
    
Metrics:
  DatabaseConnections:
    description: "Database Connections Count"
    type: gauge
  CPUUtilization:
    description: "Database CPUUtilization"
    type: gauge
  FreeableMemory:
    description: "Database FreeableMemory"
    type: gauge
  DMLLatency:
    description: "Database DMLLatency"
    type: gauge
  CommitLatency:
    description: "Database CommitLatency"
    type: gauge
```
- AWS.period : Metric Collect 주기
- AWS.logmode : debug 인경우 Terminal에 수집 메트릭 데이터가 표시

- Metrics: 수집해야하는 Metric을 정의 합니다. 

## Metric
---
-  ActiveTransactions
-  AuroraBinlogReplicaLag
-  AuroraGlobalDBDataTransferBytes
-  AuroraGlobalDBReplicatedWriteIO
-  AuroraGlobalDBReplicationLag
-  AuroraReplicaLag
-  AuroraReplicaLagMaximum
-  AuroraReplicaLagMinimum
-  AvailabilityPercentage
-  BacktrackChangeRecordsCreationRate
-  BacktrackChangeRecordsStored
-  BacktrackWindowActual
-  BacktrackWindowAlert
-  BackupRetentionPeriodStorageUsed
-  BinLogDiskUsage
-  BlockedTransactions
-  BufferCacheHitRatio
-  BurstBalance
-  CPUCreditBalance
-  CPUCreditUsage
-  CPUUtilization
-  ClientConnections
-  ClientConnectionsClosed
-  ClientConnectionsNoTLS
-  ClientConnectionsReceived
-  ClientConnectionsSetupFailedAuth
-  ClientConnectionsSetupSucceeded
-  ClientConnectionsTLS
-  CommitLatency
-  CommitThroughput
-  DDLLatency
-  DDLThroughput
-  DMLLatency
-  DMLThroughput
-  DatabaseConnectionRequests
-  DatabaseConnectionRequestsWithTLS
-  DatabaseConnections
-  DatabaseConnectionsBorrowLatency
-  DatabaseConnectionsCurrentlyBorrowed
-  DatabaseConnectionsCurrentlyInTransaction
-  DatabaseConnectionsCurrentlySessionPinned
-  DatabaseConnectionsSetupFailed
-  DatabaseConnectionsSetupSucceeded
-  DatabaseConnectionsWithTLS
-  Deadlocks
-  DeleteLatency
-  DeleteThroughput
-  DiskQueueDepth
-  EngineUptime
-  FailedSQLServerAgentJobsCount
-  FreeLocalStorage
-  FreeStorageSpace
-  FreeableMemory
-  InsertLatency
-  InsertThroughput
-  LoginFailures
-  MaxDatabaseConnectionsAllowed
-  MaximumUsedTransactionIDs
-  NetworkReceiveThroughput
-  NetworkThroughput
-  NetworkTransmitThroughput
-  OldestReplicationSlotLag
-  Queries
-  QueryDatabaseResponseLatency
-  QueryRequests
-  QueryRequestsNoTLS
-  QueryRequestsTLS
-  QueryResponseLatency
-  RDSToAuroraPostgreSQLReplicaLag
-  ReadIOPS
-  ReadLatency
-  ReadThroughput
-  ReplicaLag
-  ReplicationSlotDiskUsage
-  ResultSetCacheHitRatio
-  SelectLatency
-  SelectThroughput
-  ServerlessDatabaseCapacity
-  SnapshotStorageUsed
-  SwapUsage
-  TotalBackupStorageBilled
-  TransactionLogsDiskUsage
-  TransactionLogsGeneration
-  UpdateLatency
-  UpdateThroughput
-  VolumeBytesUsed
-  VolumeReadIOPs
-  VolumeWriteIOPs
-  WriteIOPS
-  WriteLatency
-  WriteThroughput

## Uasge
---
```
./cloud_exporter -conf=./configure.yml -exporter.port=9104 -region=ap-northeast-2 -instance=test-cluster-prod-1
```
- conf : Configure 파일
- exporter.port : Exporter 가 사용할 포트
- region : AWS 메트릭 수집 리전
- instance : RDS Instance Identifier