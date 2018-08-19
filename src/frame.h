#ifndef FRAME_H_
#define FRAME_H_

#include <cstdint>
#include <string>
#include <vector>

struct frame
{
    using data_type = std::vector<std::vector<uint8_t>>;

    frame(std::string str, size_t minimum_width);
    virtual ~frame() = default;

    size_t get_width() const { return this->width; }
    std::vector<uint8_t> concatenated_lines() const;
    frame sub_frame(int start, int end) const;

  private:
    data_type data;
    size_t width;

    void modify(int start, int end, const data_type &lines);
    void text(int pos, std::u16string str);
};

#endif
