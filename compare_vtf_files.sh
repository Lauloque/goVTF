#!/bin/bash

# Compare first 128 bytes (header + some data)
diff <(hexdump -C out/test_512.vtf | head -n 8) \
     <(hexdump -C controlTest/controlTest_512_DXT1-slowQ.vtf | head -n 8)

# Check specific fields
echo "=== Our file ==="
hexdump -C out/test_512.vtf -s 52 -n 4  # HighResFormat
hexdump -C out/test_512.vtf -s 56 -n 1  # MipmapCount
hexdump -C out/test_512.vtf -s 61 -n 2  # LowRes dimensions
hexdump -C out/test_512.vtf -s 68 -n 4  # NumResources

echo "=== Control file ==="
hexdump -C controlTest/controlTest_512_DXT1-slowQ.vtf -s 52 -n 4
hexdump -C controlTest/controlTest_512_DXT1-slowQ.vtf -s 56 -n 1
hexdump -C controlTest/controlTest_512_DXT1-slowQ.vtf -s 61 -n 2
hexdump -C controlTest/controlTest_512_DXT1-slowQ.vtf -s 68 -n 4
