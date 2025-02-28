# RedPebble
- # Roadmap version 0.1
  ## Network
  - [ ] Client connection / request / respond
  - [ ] RESP protocol
  - [ ] Design handle multiple threadings
  - [ ] Partioned data based on key to each threads
  - [ ] Transaction between threads
  - [ ] Handling concurrency
  ## Storage Engine
  - [ ] Able to register commands
  - [ ] Embeedeing repebble into main source
  - [ ] Flushing data from memory to disk
  - [ ] Implementing basic command GET, SET and TTL
  - [ ] Active key expire
  ##  Data structures
  - [ ] Implementing data structures
    - [ ] String
    - [ ] List
    - [ ] Set
    - [ ] Sorted Set
    - [ ] Hash
    - [ ] ...
  ## Eviction strategies
    - [ ] Implementing LRU algorithms
  ## Signals
    - [ ] Graceful shutdown
  ## Advanced Datastrctures
    - [ ] Skiplist
    - [ ] Sotred set
    - [ ] ....
  ## Distributed system
    - [ ] Cluster failure handling
    - [ ] Replication
    - [ ] Sharding data
    - [ ] Consistency and HA checks
  ## Backup/Restore
    - [ ] AOF and RBD files for back up and restore
  ## Benchmarking
  - [ ] Tests
    - [ ] For existing commands
    - [ ] For key expirer
  - [ ] Alpha Release

- [ ] TODO beside Roadmap
  - [ ] Persistence
  - [ ] Redis config
    - [ ] Default Redis config format
    - [ ] YAML support
    - [ ] JSON support
  - [ ] Pub/Sub
  - [ ] Redis modules
  - [ ] Benchmarks
  - [ ] Master-slave replication
  - [ ] Cluster support
  - [ ] ...
