from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone
import os

from app.network_impairment import should_drop
from app.orbital_mechanics import compute_window


@dataclass(slots=True)
class LinkQueue:
    rover_to_earth: list[dict] = field(default_factory=list)
    earth_to_rover: list[dict] = field(default_factory=list)
    one_way_delay_seconds: int = field(default_factory=lambda: int(os.getenv("SIM_DELAY_SECONDS", "0")))
    packet_loss_rate: float = field(default_factory=lambda: float(os.getenv("SIM_PACKET_LOSS_RATE", "0.0")))
    cycle_seconds: int = field(default_factory=lambda: int(os.getenv("SIM_WINDOW_CYCLE_SECONDS", "45")))
    active_seconds: int = field(default_factory=lambda: int(os.getenv("SIM_WINDOW_ACTIVE_SECONDS", "25")))

    def window(self) -> dict[str, object]:
        window = compute_window(datetime.now(timezone.utc), self.cycle_seconds, self.active_seconds)
        window["one_way_delay_seconds"] = self.one_way_delay_seconds
        window["packet_loss_rate"] = self.packet_loss_rate
        return window

    def enqueue(self, direction: str, node_id: str, bundles: list[dict]) -> list[str]:
        now = datetime.now(timezone.utc)
        accepted: list[str] = []
        for bundle in bundles:
            bundle_id = bundle.get("id", "")
            if should_drop(bundle_id, self.packet_loss_rate):
                continue
            entry = {
                "node_id": node_id,
                "available_at": (now + timedelta(seconds=self.one_way_delay_seconds)).isoformat().replace("+00:00", "Z"),
                "bundle": bundle,
            }
            if direction == "downlink":
                self.rover_to_earth.append(entry)
            else:
                self.earth_to_rover.append(entry)
            accepted.append(bundle_id)
        return accepted

    def poll(self, direction: str, node_id: str) -> list[dict]:
        source = self.rover_to_earth if direction == "earth" else self.earth_to_rover
        now = datetime.now(timezone.utc)
        delivered: list[dict] = []
        retained: list[dict] = []
        for entry in source:
            available_at = datetime.fromisoformat(entry["available_at"].replace("Z", "+00:00"))
            if entry["node_id"] == node_id and available_at <= now:
                delivered.append(entry["bundle"])
            else:
                retained.append(entry)
        if direction == "earth":
            self.rover_to_earth = retained
        else:
            self.earth_to_rover = retained
        return delivered
