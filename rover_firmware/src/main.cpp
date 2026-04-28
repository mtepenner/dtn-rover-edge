#include <Arduino.h>

namespace ultrasonic { float sampleClearanceMeters(unsigned long tick); }
namespace imu {
float sampleTiltDegrees(unsigned long tick);
float sampleHeadingDegrees(unsigned long tick);
}
namespace motors {
struct DriveCommand {
    float left;
    float right;
    bool emergency_stop;
};
DriveCommand computeDriveCommand(float clearance_meters, float tilt_degrees);
}
namespace serial_bridge {
void publishTelemetry(float x_meters, float y_meters, float heading_degrees, float clearance_meters, float tilt_degrees, bool hazard_stop);
bool tryReadCommand(float& target_x, float& target_y);
}

namespace {
float rover_x = 0.0f;
float rover_y = 0.0f;
float target_x = 8.0f;
float target_y = 2.0f;
}

void setup() {
    Serial.begin(115200);
}

void loop() {
    const unsigned long tick = millis();
    const float clearance = ultrasonic::sampleClearanceMeters(tick);
    const float tilt = imu::sampleTiltDegrees(tick);
    const float heading = imu::sampleHeadingDegrees(tick);
    motors::DriveCommand command = motors::computeDriveCommand(clearance, tilt);

    if (serial_bridge::tryReadCommand(target_x, target_y)) {
        // Target updates arrive asynchronously from the companion computer.
    }

    const float dx = target_x - rover_x;
    const float dy = target_y - rover_y;
    const float step_scale = command.emergency_stop ? 0.0f : 0.03f;
    rover_x += dx * step_scale;
    rover_y += dy * step_scale;

    serial_bridge::publishTelemetry(rover_x, rover_y, heading, clearance, tilt, command.emergency_stop);
    delay(250);
}
