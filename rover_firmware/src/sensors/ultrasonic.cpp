#include <cmath>

namespace ultrasonic {

float sampleClearanceMeters(unsigned long tick) {
    const float phase = static_cast<float>(tick) / 900.0f;
    return 0.55f + 0.22f * std::sin(phase) + 0.09f * std::cos(phase * 0.45f);
}

}  // namespace ultrasonic
