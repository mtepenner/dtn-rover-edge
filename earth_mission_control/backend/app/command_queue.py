from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime, timezone


@dataclass(slots=True)
class CommandQueue:
    sequence: int = 0
    queued: list[dict] = field(default_factory=list)
    sent: list[dict] = field(default_factory=list)

    def create_command(self, action: str, waypoint_m: list[float], speed_mps: float) -> dict:
        self.sequence += 1
        command = {
            "id": f"earth-cmd-{self.sequence:04d}",
            "action": action,
            "waypoint_m": waypoint_m,
            "parameters": {"speed_mps": speed_mps},
            "created_at": datetime.now(timezone.utc).isoformat().replace("+00:00", "Z"),
        }
        self.queued.append(command)
        return command

    def mark_sent(self, command: dict) -> None:
        self.sent.append(command)
        self.queued = [item for item in self.queued if item["id"] != command["id"]]
        self.sent = self.sent[-80:]
