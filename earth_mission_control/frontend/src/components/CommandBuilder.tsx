import { FormEvent, useState } from 'react';

import { CommandRecord } from '../types';

type Props = {
  onSubmit: (payload: { action: string; waypointX: number; waypointY: number; speed: number }) => Promise<void>;
  sentCommands: CommandRecord[];
};

export function CommandBuilder({ onSubmit, sentCommands }: Props) {
  const [action, setAction] = useState('navigate');
  const [waypointX, setWaypointX] = useState(16);
  const [waypointY, setWaypointY] = useState(4);
  const [speed, setSpeed] = useState(0.35);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    await onSubmit({ action, waypointX, waypointY, speed });
  };

  return (
    <section className="panel stack-gap">
      <div className="section-heading">
        <span>Command Builder</span>
        <small>{sentCommands.length} sent commands</small>
      </div>
      <form className="command-form" onSubmit={submit}>
        <label>
          <span>Action</span>
          <select value={action} onChange={(event) => setAction(event.target.value)}>
            <option value="navigate">Navigate</option>
            <option value="sample">Collect sample</option>
            <option value="survey">Survey ridge</option>
          </select>
        </label>
        <label>
          <span>Waypoint X</span>
          <input type="number" value={waypointX} onChange={(event) => setWaypointX(Number(event.target.value))} />
        </label>
        <label>
          <span>Waypoint Y</span>
          <input type="number" value={waypointY} onChange={(event) => setWaypointY(Number(event.target.value))} />
        </label>
        <label>
          <span>Speed (m/s)</span>
          <input type="number" step="0.05" value={speed} onChange={(event) => setSpeed(Number(event.target.value))} />
        </label>
        <button type="submit" className="primary-button">Queue DTN Command</button>
      </form>
      <div className="command-history">
        {sentCommands.slice(-4).reverse().map((command) => (
          <article key={command.id} className="history-card">
            <strong>{command.action}</strong>
            <span>{command.waypoint_m[0].toFixed(1)}, {command.waypoint_m[1].toFixed(1)} m</span>
          </article>
        ))}
      </div>
    </section>
  );
}
