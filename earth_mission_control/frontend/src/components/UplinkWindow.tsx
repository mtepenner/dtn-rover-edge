import { WindowStatus } from '../types';

type Props = {
  windowStatus: WindowStatus | null;
};

export function UplinkWindow({ windowStatus }: Props) {
  if (!windowStatus) {
    return (
      <section className="panel stack-gap">
        <div className="section-heading"><span>Uplink Window</span><small>Awaiting link status</small></div>
      </section>
    );
  }

  const nextWindow = new Date(windowStatus.next_window_start);
  const currentWindowEnd = new Date(windowStatus.current_window_end);
  const currentTime = new Date(windowStatus.current_time);
  const countdownSeconds = Math.max(0, Math.round((windowStatus.active ? currentWindowEnd.getTime() - currentTime.getTime() : nextWindow.getTime() - currentTime.getTime()) / 1000));

  return (
    <section className="panel stack-gap">
      <div className="section-heading">
        <span>Uplink Window</span>
        <small>{windowStatus.active ? 'Orbiter link active' : 'Waiting for line of sight'}</small>
      </div>
      <div className={`window-badge ${windowStatus.active ? 'active' : ''}`}>{windowStatus.active ? 'ACTIVE' : 'STANDBY'}</div>
      <div className="metric-row"><span>Countdown</span><strong>{countdownSeconds}s</strong></div>
      <div className="metric-row"><span>One-way delay</span><strong>{windowStatus.one_way_delay_seconds}s</strong></div>
      <div className="metric-row"><span>Packet loss</span><strong>{(windowStatus.packet_loss_rate * 100).toFixed(1)}%</strong></div>
    </section>
  );
}
