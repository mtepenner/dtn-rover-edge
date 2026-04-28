# Delay-Tolerant Network (DTN) Rover Edge Compute

A robust, emulated Delay-Tolerant Networking (DTN) architecture designed for autonomous planetary rovers. This project demonstrates high-latency, unreliable communication handling via the Bundle Protocol (RFC 5050/9171).

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Technologies](#technologies)
- [Installation](#installation)
- [License](#license)

## 🚀 Features
- **DTN Bundle Protocol**: Implements RFC 5050/9171 for reliable data delivery over unstable links.
- **Deep Space Simulation**: Models massive latency (2-hour delay) and atmospheric noise/radiation.
- **Autonomous Navigation**: ESP32-based collision avoidance and tilt-detection.
- **Asynchronous Mission Control**: Earth-side UI for scheduling commands and visualizing delayed telemetry.

## 🏗️ Architecture
The system is divided into four primary components:
1.  **Rover Firmware (C++)**: Handles real-time hardware survival reflexes.
2.  **Edge Daemon (Go)**: The onboard brain, managing data bundles and protocol adherence.
3.  **Deep Space Link (Python)**: Simulates the physics and network impairments of space communications.
4.  **Mission Control (React/Go)**: The asynchronous interface for commanding the rover.

## 🛠️ Technologies
- **Embedded**: C++ (ESP32/FreeRTOS)
- **Daemon**: Go (Bundle Protocol)
- **Simulation/Backend**: Python
- **Frontend**: React (TypeScript)
- **Infrastructure**: Docker & Docker Compose

## 📥 Installation
1. Clone the repository: `git clone https://github.com/mtepenner/dtn-rover-edge.git`
2. Launch the full environment: `docker-compose up`

For firmware-only builds, if `pio` is not available on your shell `PATH`, run PlatformIO through Python instead: `python -m platformio run -d rover_firmware`

## ⚖️ License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
