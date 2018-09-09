#include "max7219.h"
#include "sign.h"

#include <iostream>

int main() {
  try {
    max7219 max("/dev/spidev0.1", 4);
    max.init();
    bus_sign([&max](const frame &f) { max.display(f); });
  } catch (const std::exception &e) {
    std::cerr << "Error: " << e.what() << '\n';
    return 1;
  }
  return 0;
}
