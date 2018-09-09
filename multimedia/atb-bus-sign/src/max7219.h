#ifndef MAX7219_H_
#define MAX7219_H_

#include "frame.h"
#include <cstdint>
#include <string>

class max7219 {
  enum class reg : uint8_t {
    no_op = 0x0,
    digit_0 = 0x1,
    digit_1 = 0x2,
    digit_2 = 0x3,
    digit_3 = 0x4,
    digit_4 = 0x5,
    digit_5 = 0x6,
    digit_6 = 0x7,
    digit_7 = 0x8,
    decode_mode = 0x9,
    intensity = 0xA,
    scan_limit = 0xB,
    shutdown = 0xC,
    display_test = 0xF
  };

  static const uint8_t spi_mode = 0;
  static const uint32_t speed_hz = 1000000;
  static const uint8_t bits_per_word = 8;

  const int cascade;
  const int spi_dev_fd;

  max7219(const max7219 &) = delete;
  max7219 &operator=(const max7219 &) = delete;

  void spi_transfer(const uint8_t *data, int length);
  void write_to_all(reg r, uint8_t data);
  void display_line(int line, const uint8_t *data);

public:
  max7219(std::string device, int cascade);
  virtual ~max7219();

  void init();
  void display(const frame &f);
};

#endif
