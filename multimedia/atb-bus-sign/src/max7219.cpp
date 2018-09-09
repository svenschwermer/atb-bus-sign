#include "max7219.h"

#include <linux/spi/spidev.h>
#include <linux/types.h>

#include <fcntl.h>
#include <sys/ioctl.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include <cstring>
#include <initializer_list>
#include <system_error>
#include <vector>

max7219::max7219(std::string device, int cascade)
    : cascade(cascade), spi_dev_fd(::open(device.c_str(), O_RDWR)) {
  if (spi_dev_fd == -1)
    throw std::system_error(errno, std::generic_category());

  auto mode = spi_mode;
  if (::ioctl(spi_dev_fd, SPI_IOC_WR_MODE, &mode) < 0)
    throw std::system_error(errno, std::generic_category());
}

max7219::~max7219() {
  if (spi_dev_fd != -1)
    ::close(spi_dev_fd);
}

void max7219::init() {
  std::initializer_list<std::pair<reg, uint8_t>> sequence{
      {reg::decode_mode, 0x00},  // No decoding
      {reg::intensity, 0x01},    // Minimal intensity
      {reg::scan_limit, 0x07},   // Scan all digits (0..7)
      {reg::display_test, 0x00}, // Disable display test
      {reg::shutdown, 0x01},     // Normal operation
  };
  for (const auto x : sequence)
    write_to_all(x.first, x.second);
}

void max7219::display(const frame &f) {
  auto bytes = f.concatenated_lines();
  if (bytes.size() != 8 * cascade)
    throw std::runtime_error(__func__);
  for (int line = 0; line < 8; ++line) {
    auto line_data = bytes.data() + line * cascade;
    display_line(line, line_data);
  }
}

void max7219::spi_transfer(const uint8_t *data, int length) {
  std::vector<uint8_t> rx_buf(length);
  struct spi_ioc_transfer xfer;
  memset(&xfer, 0, sizeof(xfer));
  xfer.tx_buf = reinterpret_cast<uint64_t>(data);
  xfer.rx_buf = reinterpret_cast<uint64_t>(rx_buf.data());
  xfer.len = length;
  xfer.speed_hz = speed_hz;
  xfer.bits_per_word = bits_per_word;

  if (::ioctl(spi_dev_fd, SPI_IOC_MESSAGE(1), &xfer) < 0)
    throw std::system_error(errno, std::generic_category());
}

void max7219::write_to_all(reg r, uint8_t data) {
  std::vector<uint8_t> buf(2 * cascade);
  for (int i = 0; i < cascade; ++i) {
    buf[2 * i] = static_cast<uint8_t>(r);
    buf[2 * i + 1] = data;
  }
  spi_transfer(buf.data(), buf.size());
}

void max7219::display_line(int line, const uint8_t *data) {
  if (line < 0 || line > 7)
    throw std::runtime_error(__func__);
  std::vector<uint8_t> buf(2 * cascade);
  for (int i = 0; i < cascade; ++i) {
    buf[2 * i] = static_cast<uint8_t>(line + 1);
    buf[2 * i + 1] = data[i];
  }
  spi_transfer(buf.data(), buf.size());
}
