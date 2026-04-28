from __future__ import annotations

import os

import httpx
from fastapi import FastAPI
from pydantic import BaseModel, Field

from app.command_queue import CommandQueue
from app.database import TelemetryDatabase
from app.dtn_receiver import unpack_telemetry_bundles


class CommandRequest(BaseModel):
    action: str = Field(default="navigate")
    waypoint_x_m: float
    waypoint_y_m: float
    speed_mps: float = Field(default=0.35, ge=0.05, le=1.5)


app = FastAPI(title="Earth Mission Control Backend")
database = TelemetryDatabase()
queue = CommandQueue()
deep_space_link_url = os.getenv("DEEP_SPACE_LINK_URL", "http://127.0.0.1:8082")
rover_node_id = os.getenv("ROVER_NODE_ID", "titan-rover-7")


@app.get("/health")
async def health() -> dict[str, object]:
    return {
        "service": "earth-mission-control-backend",
        "status": "ok",
        "telemetry_points": len(database.telemetry_history),
        "queued_commands": len(queue.queued),
    }


@app.get("/telemetry")
async def telemetry() -> list[dict]:
    return database.telemetry_history


@app.get("/commands")
async def commands() -> dict[str, list[dict]]:
    return {"queued": queue.queued, "sent": queue.sent}


@app.get("/window")
async def window() -> dict:
    async with httpx.AsyncClient(timeout=4.0) as client:
        response = await client.get(f"{deep_space_link_url}/window")
        response.raise_for_status()
        return response.json()


@app.post("/commands")
async def create_command(request: CommandRequest) -> dict:
    command = queue.create_command(request.action, [request.waypoint_x_m, request.waypoint_y_m], request.speed_mps)
    bundle = {
        "id": command["id"],
        "sequence": queue.sequence,
        "source": "earth-mission-control",
        "destination": rover_node_id,
        "type": "command",
        "priority": "expedited",
        "created_at": command["created_at"],
        "ttl_minutes": 720,
        "payload": command,
    }

    async with httpx.AsyncClient(timeout=4.0) as client:
        response = await client.post(f"{deep_space_link_url}/relay/uplink", json={"node_id": rover_node_id, "bundles": [bundle]})
        response.raise_for_status()
    queue.mark_sent(command)
    return command


@app.post("/sync")
async def sync_downlink() -> dict[str, object]:
    async with httpx.AsyncClient(timeout=4.0) as client:
        response = await client.get(f"{deep_space_link_url}/poll/earth", params={"node_id": rover_node_id})
        response.raise_for_status()
        bundles = response.json().get("bundles", [])

    telemetry_payloads = unpack_telemetry_bundles(bundles)
    for telemetry in telemetry_payloads:
        database.ingest(telemetry)

    return {"received_bundles": len(bundles), "received_telemetry": len(telemetry_payloads)}
