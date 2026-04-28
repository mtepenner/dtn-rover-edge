#include <cmath>

namespace imu {

float sampleTiltDegrees(unsigned long tick) {
    const float phase = static_cast<float>(tick) / 1100.0f;
    return 4.2f * std::sin(phase * 0.8f) + 1.1f * std::cos(phase * 0.23f);
}

float sampleHeadingDegrees(unsigned long tick) {
    const float phase = static_cast<float>(tick) / 1500.0f;
    const float wrapped = std::fmod(phase * 37.0f, 360.0f);
    return wrapped < 0.0f ? wrapped + 360.0f : wrapped;
}

}  // namespace imu