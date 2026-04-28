from __future__ import annotations

from datetime import datetime, timedelta, timezone


def compute_window(now: datetime, cycle_seconds: int, active_seconds: int) -> dict[str, object]:
    cycle_offset = int(now.timestamp()) % cycle_seconds
    current_cycle_start = now - timedelta(seconds=cycle_offset)
    active = cycle_offset < active_seconds
    if active:
        current_window_end = current_cycle_start + timedelta(seconds=active_seconds)
        next_window_start = current_cycle_start + timedelta(seconds=cycle_seconds)
    else:
        current_window_end = current_cycle_start
        next_window_start = current_cycle_start + timedelta(seconds=cycle_seconds)

    return {
        "active": active,
        "current_time": now.astimezone(timezone.utc).isoformat().replace("+00:00", "Z"),
        "current_window_end": current_window_end.astimezone(timezone.utc).isoformat().replace("+00:00", "Z"),
        "next_window_start": next_window_start.astimezone(timezone.utc).isoformat().replace("+00:00", "Z"),
    }
