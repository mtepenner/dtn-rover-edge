import { useEffect, useState } from 'react';

import { AsynchronousMap } from './components/AsynchronousMap';
import { CommandBuilder } from './components/CommandBuilder';
import { UplinkWindow } from './components/UplinkWindow';
import { CommandRecord, Telemetry, WindowStatus } from './types';

const backendUrl = import.meta.env.VITE_EARTH_BACKEND_URL ?? 'http://127.0.0.1:8083';

export default function App() {
  const [telemetry, setTelemetry] = useState<Telemetry[]>([]);
  const [windowStatus, setWindowStatus] = useState<WindowStatus | null>(null);
  const [sentCommands, setSentCommands] = useState<CommandRecord[]>([]);
  const [status, setStatus] = useState('Syncing with delayed network');

  useEffect(() => {
    let cancelled = false;

    const poll = async () => {
      try {
        await fetch(`${backendUrl}/sync`, { method: 'POST' });
        const [telemetryResponse, windowResponse, commandsResponse] = await Promise.all([
          fetch(`${backendUrl}/telemetry`),
          fetch(`${backendUrl}/window`),
          fetch(`${backendUrl}/commands`),
        ]);
        if (!telemetryResponse.ok || !windowResponse.ok || !commandsResponse.ok) {
          throw new Error('Mission control endpoints unavailable');
        }

        const nextTelemetry = (await telemetryResponse.json()) as Telemetry[];
        const nextWindow = (await windowResponse.json()) as WindowStatus;
        const nextCommands = (await commandsResponse.json()) as { queued: CommandRecord[]; sent: CommandRecord[] };

        if (!cancelled) {
          setTelemetry(nextTelemetry);
          setWindowStatus(nextWindow);
          setSentCommands(nextCommands.sent);
          setStatus(`Last sync ${new Date().toLocaleTimeString()}`);
        }
      } catch (error) {
        if (!cancelled) {
          setStatus(error instanceof Error ? error.message : 'Mission control sync failed');
        }
      }
    };

    poll();
    const intervalId = window.setInterval(poll, 1500);
    return () => {
      cancelled = true;
      window.clearInterval(intervalId);
    };
  }, []);

  const submitCommand = async (payload: { action: string; waypointX: number; waypointY: number; speed: number }) => {
    const response = await fetch(`${backendUrl}/commands`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: payload.action,
        waypoint_x_m: payload.waypointX,
        waypoint_y_m: payload.waypointY,
        speed_mps: payload.speed,
      }),
    });
    const command = (await response.json()) as CommandRecord;
    setSentCommands((current) => [...current, command]);
  };

  return (
    <main className="shell">
      <section className="hero panel">
        <div>
          <p className="eyebrow">Earth Mission Control</p>
          <h1>Asynchronous command and telemetry console for a rover operating beyond continuous contact.</h1>
        </div>
        <div className="hero-meta">
          <span>{status}</span>
          <span>{telemetry.length ? `Latest fix ${new Date(telemetry[telemetry.length - 1].timestamp).toLocaleTimeString()}` : 'No downlink yet'}</span>
        </div>
      </section>

      <section className="dashboard-grid">
        <AsynchronousMap telemetry={telemetry} />
        <div className="sidebar-stack">
          <UplinkWindow windowStatus={windowStatus} />
          <CommandBuilder onSubmit={submitCommand} sentCommands={sentCommands} />
        </div>
      </section>
    </main>
  );
}
