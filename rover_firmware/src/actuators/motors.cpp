#include <algorithm>

namespace motors {

struct DriveCommand {
    float left;
    float right;
    bool emergency_stop;
};

DriveCommand computeDriveCommand(float clearance_meters, float tilt_degrees) {
    DriveCommand command{};
    command.emergency_stop = clearance_meters < 0.28f || tilt_degrees > 11.0f || tilt_degrees < -11.0f;
    if (command.emergency_stop) {
        return command;
    }

    const float speed_scale = std::clamp((clearance_meters - 0.25f) / 0.55f, 0.0f, 1.0f);
    const float tilt_bias = std::clamp(tilt_degrees / 15.0f, -0.35f, 0.35f);
    command.left = 0.45f * speed_scale - tilt_bias;
    command.right = 0.45f * speed_scale + tilt_bias;
    return command;
}

}  // namespace motors
