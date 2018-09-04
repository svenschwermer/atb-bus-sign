#include "frame.h"
#include <iostream>

int main()
{
    auto f = frame("waiting...  ", 32);
    for (int i = 0; i < 3; ++i)
        std::cout << f.sub_frame(i, i + 32);
    return 0;
}