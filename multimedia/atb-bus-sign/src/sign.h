#ifndef SIGN_H_
#define SIGN_H_

#include "atb.h"
#include "frame.h"

#include <chrono>
#include <future>
#include <string>
#include <thread>

template <typename DisplayFunc> static void bus_sign(DisplayFunc display) {
  using namespace std::chrono_literals;
  using namespace std::string_literals;

  std::vector<std::string> lines = {"11"};
  auto poll = [lines](std::chrono::milliseconds pre_wait) {
    return std::async(std::launch::async, get_bus_sign, "16011297"s, lines,
                      pre_wait);
  };
  auto future = poll(0s);

  auto f = frame("waiting...  ", 32);
  while (true) {
    if (future.wait_for(0s) == std::future_status::ready) {
      auto str = future.get() + " | ";
      f = frame(str, 32);
      future = poll(10s);
    }
    for (int i = 0; i < f.get_width(); ++i) {
      auto sub = f.sub_frame(i, i + 32);
      display(sub);
      std::this_thread::sleep_for(50ms);
    }
  }
}

#endif
