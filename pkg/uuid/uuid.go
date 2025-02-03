package uuid

type Generator interface {
	// NewV1 generates a UUID based on the current timestamp and the MAC address of the machine.
	// Ensures uniqueness across space and time but may expose the MAC address.
	NewV1() (string, error)

	// NewV4 Generates a UUID using random or pseudo-random numbers.
	// This is the most commonly used version and provides a good balance between uniqueness and privacy.
	NewV4() (string, error)

	// NewV6 generates an ordered time-based UUID.
	// Similar to Version 1 but with a reordered timestamp to improve locality and sorting efficiency.
	// Useful for databases where ordered UUIDs are beneficial.
	NewV6() (string, error)

	// NewV7 generates a UUID based on a custom epoch (e.g., Unix epoch)
	// and includes a timestamp in the most significant bits.
	// This version is designed for better sorting and indexing in distributed systems.
	NewV7() (string, error)
}
