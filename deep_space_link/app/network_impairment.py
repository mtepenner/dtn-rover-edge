from __future__ import annotations

import hashlib


def should_drop(bundle_id: str, loss_rate: float) -> bool:
    if loss_rate <= 0.0:
        return False
    digest = hashlib.sha256(bundle_id.encode("utf-8")).digest()
    bucket = int.from_bytes(digest[:2], "big") / 65535.0
    return bucket < loss_rate
