#include <Arduino.h>

namespace serial_bridge {

void publishTelemetry(
    float x_meters,
    float y_meters,
    float heading_degrees,
    float clearance_meters,
    float tilt_degrees,
    bool hazard_stop
) {
    Serial.print("{\"x_m\":");
    Serial.print(x_meters, 3);
    Serial.print(",\"y_m\":");
    Serial.print(y_meters, 3);
    Serial.print(",\"heading_deg\":");
    Serial.print(heading_degrees, 2);
    Serial.print(",\"clearance_m\":");
    Serial.print(clearance_meters, 3);
    Serial.print(",\"tilt_deg\":");
    Serial.print(tilt_degrees, 2);
    Serial.print(",\"hazard_stop\":");
    Serial.print(hazard_stop ? "true" : "false");
    Serial.println("}");
}

bool tryReadCommand(float& target_x, float& target_y) {
    if (!Serial.available()) {
        return false;
    }

    String command = Serial.readStringUntil('\n');
    command.trim();
    const int separator = command.indexOf(',');
    if (separator < 0) {
        return false;
    }

    target_x = command.substring(0, separator).toFloat();
    target_y = command.substring(separator + 1).toFloat();
    return true;
}

}  // namespace serial_bridge
