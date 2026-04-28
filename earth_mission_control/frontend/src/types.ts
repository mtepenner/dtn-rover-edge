export type Telemetry = {
  timestamp: string;
  position_m: [number, number];
  heading_deg: number;
  clearance_m: number;
  tilt_deg: number;
  battery_pct: number;
  hazard_stop: boolean;
  mode: string;
};

export type WindowStatus = {
  active: boolean;
  current_time: string;
  current_window_end: string;
  next_window_start: string;
  one_way_delay_seconds: number;
  packet_loss_rate: number;
};

export type CommandRecord = {
  id: string;
  action: string;
  waypoint_m: [number, number];
  parameters: Record<string, number>;
  created_at: string;
};
