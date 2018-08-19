#include "max7219.h"
#include <chrono>
#include <iostream>
#include <thread>

static void display(max7219 &max);

int main()
{
    try
    {
        max7219 max("/dev/spidev0.1", 4);
        max.init();
        display(max);
    }
    catch (const std::exception &e)
    {
        std::cerr << "Error: " << e.what() << '\n';
        return 1;
    }
    return 0;
}

static void display(max7219 &max)
{
    auto f = frame("waiting...  ", 32);
    while (true)
    {
        for (int i = 0; i < f.get_width(); ++i)
        {
			auto sub = f.sub_frame(i, i+32);
			max.display(sub);
            std::this_thread::sleep_for(std::chrono::milliseconds(50));
		}
    }
}
