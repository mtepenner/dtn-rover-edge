from __future__ import annotations

from dataclasses import dataclass, field


@dataclass(slots=True)
class TelemetryDatabase:
    telemetry_history: list[dict] = field(default_factory=list)

    def ingest(self, telemetry: dict) -> None:
        self.telemetry_history.append(telemetry)
        self.telemetry_history = self.telemetry_history[-120:]
