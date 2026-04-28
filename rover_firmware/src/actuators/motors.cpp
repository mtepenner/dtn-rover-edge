namespace motors {

namespace {

float clampFloat(float value, float min_value, float max_value) {
    if (value < min_value) {
        return min_value;
    }
    if (value > max_value) {
        return max_value;
    }
    return value;
}

}  // namespace

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

    const float speed_scale = clampFloat((clearance_meters - 0.25f) / 0.55f, 0.0f, 1.0f);
    const float tilt_bias = clampFloat(tilt_degrees / 15.0f, -0.35f, 0.35f);
    command.left = 0.45f * speed_scale - tilt_bias;
    command.right = 0.45f * speed_scale + tilt_bias;
    return command;
}

}  // namespace motors
