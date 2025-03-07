# Vector database that compatible with SQL protocol
## Query patterns
Create: CREATE TABLE items (id bigserial PRIMARY KEY, embedding vector(3));
Insert: INSERT INTO items (embedding) VALUES ('[1,2,3]'), ('[4,5,6]');
Search: SELECT * FROM items ORDER BY embedding <-> '[3,1,2]' LIMIT 5;
Querying:
Get the nearest neighbors to a vector

SELECT * FROM items ORDER BY embedding <-> '[3,1,2]' LIMIT 5;
Supported distance functions are:

<-> - L2 distance
<#> - (negative) inner product
<=> - cosine distance
<+> - L1 distance
<~> - Hamming distance (binary vectors)
<%> - Jaccard distance (binary vectors)
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
