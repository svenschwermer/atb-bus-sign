#include "atb.h"
#include "frame.h"
#include "max7219.h"

#include <chrono>
#include <future>
#include <iostream>
#include <string>
#include <thread>

static void bus_sign(max7219 &max);

int main()
{
    try
    {
        max7219 max("/dev/spidev0.1", 4);
        max.init();
        bus_sign(max);
    }
    catch (const std::exception &e)
    {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }
    return 0;
}

static void bus_sign(max7219 &max)
{
    using namespace std::chrono_literals;
    using namespace std::string_literals;

    auto poll = [](std::chrono::milliseconds pre_wait) {
        auto lines = {"11"s, "19"s};
        return std::async(std::launch::async, get_bus_sign, "16011297"s, lines, pre_wait);
    };
    auto future = poll(0s);

    auto f = frame("waiting...  ", 32);
    while (true)
    {
        if (future.valid())
        {
            auto str = future.get() + " | ";
            f = frame(str, 32);
            future = poll(10s);
        }
        for (int i = 0; i < f.get_width(); ++i)
        {
            auto sub = f.sub_frame(i, i + 32);
            max.display(sub);
            std::this_thread::sleep_for(std::chrono::milliseconds(50));
        }
    }
}
