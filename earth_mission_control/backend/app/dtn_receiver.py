from __future__ import annotations

import base64
import json


def unpack_telemetry_bundles(bundles: list[dict]) -> list[dict]:
    telemetry = []
    for bundle in bundles:
        if bundle.get("type") != "telemetry":
            continue
        payload = bundle.get("payload")
        if isinstance(payload, dict):
            telemetry.append(payload)
            continue
        if isinstance(payload, str):
            try:
                telemetry.append(json.loads(payload))
            except json.JSONDecodeError:
                telemetry.append(json.loads(base64.b64decode(payload).decode("utf-8")))
            continue
        telemetry.append(json.loads(bytes(payload).decode("utf-8")))
    return telemetry
