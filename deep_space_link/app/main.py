from __future__ import annotations

from fastapi import FastAPI
from pydantic import BaseModel

from app.simulator import LinkQueue


class RelayPayload(BaseModel):
    node_id: str
    bundles: list[dict]


app = FastAPI(title="Deep Space Link Simulator")
queue = LinkQueue()


@app.get("/health")
async def health() -> dict[str, object]:
    return {
        "service": "deep-space-link",
        "status": "ok",
        "queued_downlink": len(queue.rover_to_earth),
        "queued_uplink": len(queue.earth_to_rover),
    }


@app.get("/window")
async def window() -> dict[str, object]:
    return queue.window()


@app.post("/relay/downlink")
async def relay_downlink(payload: RelayPayload) -> dict[str, object]:
    return {"accepted_bundle_ids": queue.enqueue("downlink", payload.node_id, payload.bundles)}


@app.post("/relay/uplink")
async def relay_uplink(payload: RelayPayload) -> dict[str, object]:
    return {"accepted_bundle_ids": queue.enqueue("uplink", payload.node_id, payload.bundles)}


@app.get("/poll/earth")
async def poll_earth(node_id: str = "titan-rover-7") -> dict[str, object]:
    return {"bundles": queue.poll("earth", node_id)}


@app.get("/poll/rover")
async def poll_rover(node_id: str) -> dict[str, object]:
    return {"bundles": queue.poll("rover", node_id)}
