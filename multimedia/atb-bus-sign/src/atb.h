#ifndef ATB_H_
#define ATB_H_

#include <chrono>
#include <string>
#include <vector>

std::string get_bus_sign(std::string node_id, std::vector<std::string> lines,
                         std::chrono::milliseconds pre_wait);

#endif
