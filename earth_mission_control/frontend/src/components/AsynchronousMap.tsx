import { Telemetry } from '../types';

type Props = {
  telemetry: Telemetry[];
};

export function AsynchronousMap({ telemetry }: Props) {
  const latest = telemetry[telemetry.length - 1] ?? null;

  return (
    <section className="panel stack-gap">
      <div className="section-heading">
        <span>Asynchronous Map</span>
        <small>{telemetry.length} delayed fixes</small>
      </div>
      <svg viewBox="0 0 600 320" className="map-grid" role="img" aria-label="Delayed rover path map">
        <rect x="0" y="0" width="600" height="320" rx="22" />
        {Array.from({ length: 11 }).map((_, index) => (
          <g key={index}>
            <line x1={index * 60} y1="0" x2={index * 60} y2="320" />
            <line x1="0" y1={index * 32} x2="600" y2={index * 32} />
          </g>
        ))}
        {telemetry.map((point, index) => {
          const x = 60 + point.position_m[0] * 28;
          const y = 260 - point.position_m[1] * 20;
          const prev = telemetry[index - 1];
          return (
            <g key={`${point.timestamp}-${index}`}>
              {prev ? (
                <line
                  x1={60 + prev.position_m[0] * 28}
                  y1={260 - prev.position_m[1] * 20}
                  x2={x}
                  y2={y}
                  className="path-line"
                />
              ) : null}
              <circle cx={x} cy={y} r={index === telemetry.length - 1 ? 6 : 3.5} className="path-point" />
            </g>
          );
        })}
        {latest ? <text x="20" y="28">Last telemetry: {new Date(latest.timestamp).toLocaleTimeString()}</text> : null}
      </svg>
    </section>
  );
}
