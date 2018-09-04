#ifndef ATB_H_
#define ATB_H_

#include <chrono>
#include <initializer_list>
#include <string>

std::string get_bus_sign(std::string node_id, std::initializer_list<std::string> lines,
                         std::chrono::milliseconds pre_wait);

#endif
