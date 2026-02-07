package log

// Config holds the configuration for the log.
type Config struct {

	/*
		What is a Segment?
		It's a portion of the log that contains a store and an index.
		The store is where the actual log records are stored, while the index is a mapping of record offsets to their positions in the store.
		The Segment struct defines the parameters for how large the store and index can grow before a new segment is created.

		The Segment struct wraps the index and store bytes to coordinate operations between them.
		When the log appends a new record, the active segment will write the data to the store
		and also add new entry in the index. For reads, the segment will also have to
		look up the entry in the index and then fetch the data from the store.

		- MaxStoreBytes: The maximum number of bytes the store can use before it rolls over to a new segment.
		- MaxIndexBytes: The maximum number of bytes the index can use before it rolls over to a new segment.
		- InitialOffset: The starting offset for the log, which is useful for resuming from a specific point after a restart or failure.
	*/
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
