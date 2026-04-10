package qdrant

// ZAP wire protocol opcodes.
// Core opcodes 0x01-0x07 match the server spec.
// Extended opcodes 0x10+ are routed by the server's ZAP handler.
const (
	// Core opcodes (spec)
	opUpsertPoints     uint16 = 0x01
	opSearchPoints     uint16 = 0x02
	opDeletePoints     uint16 = 0x03
	opGetPoints        uint16 = 0x04
	opCreateCollection uint16 = 0x05
	opListCollections  uint16 = 0x06
	opHealthCheck      uint16 = 0x07

	// Points extended
	opScrollPoints          uint16 = 0x10
	opUpdateVectors         uint16 = 0x11
	opDeleteVectors         uint16 = 0x12
	opSetPayload            uint16 = 0x13
	opOverwritePayload      uint16 = 0x14
	opDeletePayload         uint16 = 0x15
	opClearPayload          uint16 = 0x16
	opCreateFieldIndex      uint16 = 0x17
	opDeleteFieldIndex      uint16 = 0x18
	opCountPoints           uint16 = 0x19
	opUpdateBatch           uint16 = 0x1A
	opQueryPoints           uint16 = 0x1B
	opQueryBatchPoints      uint16 = 0x1C
	opQueryGroupPoints      uint16 = 0x1D
	opFacetCounts           uint16 = 0x1E
	opSearchMatrixPairs     uint16 = 0x1F
	opSearchMatrixOffsets   uint16 = 0x20

	// Collections extended
	opGetCollection            uint16 = 0x30
	opDeleteCollection         uint16 = 0x31
	opUpdateCollection         uint16 = 0x32
	opCollectionExists         uint16 = 0x33
	opUpdateAliases            uint16 = 0x34
	opListCollectionAliases    uint16 = 0x35
	opListAliases              uint16 = 0x36
	opCreateShardKey           uint16 = 0x37
	opDeleteShardKey           uint16 = 0x38
	opCollectionClusterInfo    uint16 = 0x39
	opUpdateCollectionCluster  uint16 = 0x3A

	// Snapshots
	opCreateSnapshot     uint16 = 0x40
	opListSnapshots      uint16 = 0x41
	opDeleteSnapshot     uint16 = 0x42
	opCreateFullSnapshot uint16 = 0x43
	opListFullSnapshots  uint16 = 0x44
	opDeleteFullSnapshot uint16 = 0x45
)
