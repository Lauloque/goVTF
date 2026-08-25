/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

// Based onVTF 7.3 defined in https://developer.valvesoftware.com/wiki/VTF_(Valve_Texture_Format)

// VTFHeader gets go-specific padding, do not use binary.Write directly on it.
type VTFHeader struct {
	Signature     [4]byte    // File signature ("VTF\0"). (or as little-endian integer, 0x00465456)
	Version       [2]uint32  // version[0].version[1] (currently 7.2).
	HeaderSize    uint32     // Size of the header struct  (16 byte aligned; currently 80 bytes) + size of the resources dictionary (7.3+).
	Width         uint16     // Width of the largest mipmap in pixels.
	Height        uint16     // Height of the largest mipmap in pixels.
	Flags         uint32     // VTF flags.
	Frames        uint16     // Number of frames, if animated (1 for no animation).
	FirstFrame    uint16     // First frame in animation (0 based). Can be -1 in environment maps older than 7.5, meaning there are 7 faces, not 6.
	Padding0      [4]byte    // reflectivity padding (16 byte alignment).
	Reflectivity  [3]float32 // reflectivity vector.
	Padding1      [4]byte    // reflectivity padding (8 byte packing).
	BumpmapScale  float32    // Bumpmap scale.
	HighResFormat int32      // High resolution image format.
	MipmapCount   uint8      // Number of mipmaps.
	LowResFormat  int32      // Low resolution image format. This value should always be assumed to be DXT1!
	LowResWidth   uint8      // Low resolution image width.
	LowResHeight  uint8      // Low resolution image height.
	Depth         uint16     // Depth of the largest mipmap in pixels. Must be a power of 2. Is 1 for a 2D texture.
	Padding2      [3]byte    // depth padding (4 byte alignment).
	NumResources  uint32     // Number of resources this vtf has. The max appears to be 32.
	Padding3      [8]byte    // Necessary on certain compilers (changed from [8]byte to [4]byte)
}

// ResourceEntry represents a resource directory entry.
type ResourceEntry struct {
	Tag    [3]byte // A three-byte "tag" that identifies what this resource is.
	Flags  uint8   // Resource entry flags. The only known flag is 0x2, which indicates that no data chunk corresponds to this resource.
	Offset uint32  // The offset of this resource's data in the file.
}
