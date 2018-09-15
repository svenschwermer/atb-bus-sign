#ifndef FONT_H_
#define FONT_H_

#include <cstdint>
#include <vector>

struct font {
  using glyph_type = std::vector<std::vector<uint8_t>>;

  static const int width = 6; // includes spacing to the right
  static const int height = 8;

  static const glyph_type &get(char32_t c);

private:
  static const int offset = 0x20;
  static const std::vector<glyph_type> charmap;
};

#endif
