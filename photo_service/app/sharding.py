from __future__ import annotations

import uuid


def shard_index(image_id: uuid.UUID) -> int:
    """
    UUID sharding rule:
    0-3 -> shard0, 4-7 -> shard1, 8-b -> shard2, c-f -> shard3.
    """
    first = image_id.hex[0].lower()
    if first in {"0", "1", "2", "3"}:
        return 0
    if first in {"4", "5", "6", "7"}:
        return 1
    if first in {"8", "9", "a", "b"}:
        return 2
    return 3

