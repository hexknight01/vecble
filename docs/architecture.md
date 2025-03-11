# REDIS + PEBBLE DATABASE
RedPebble is the database servered as Redis on memory.

# How RedPebble multiple-threads to server the data

# How RedPebble guarantees transaction isolation between threads.

# Model HNSW graph with key-value store
Modeling Hierarchical Navigable Small World (HNSW) in a key-value database involves designing a structure that efficiently stores and retrieves graph-based data for approximate nearest neighbor (ANN) search.

Key	Value
layer1:node1	{vector: [0.1, 0.2], neighbors: [node2, node3]}
layer1:node2	{vector: [0.3, 0.4], neighbors: [node1, node4]}
layer2:node5	{vector: [0.5, 0.6], neighbors: [node6]}
