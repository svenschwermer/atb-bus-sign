#include "frame.h"
#include "font.h"

#include <algorithm>
#include <iterator>
#include <locale>
#include <codecvt>

frame::frame(std::string str, size_t minimum_width)
    : data(font::height)
{
    auto utf16 = std::wstring_convert<std::codecvt_utf8_utf16<char16_t>, char16_t>{}.from_bytes(str.data());
    width = std::max(utf16.length() * font::width, minimum_width);
    for (auto &line : data)
        line.resize(width);
    text(0, utf16);
}

std::vector<uint8_t> frame::concatenated_lines() const
{
    std::vector<uint8_t> buf;
    for (auto &line : data)
        std::copy(line.begin(), line.end(), std::back_inserter(buf));
    return buf;
}

frame frame::sub_frame(int start, int end) const
{
    frame sub(*this);
    sub.width = end - start;
    auto shift = start % 8;
    start /= 8;
    uint8_t mask = (1 << shift) - 1;
    for (auto &line : sub.data)
    {
        line.erase(line.begin(), line.begin() + start); // drop whole leading bytes
        if (shift > 0)
        {
            int j = 0;
            for (auto &b : line)
            {
                b <<= shift;
                if (j > 0)
                {
                    auto overflow = (b >> (8 - shift)) & mask;
                    line[j - 1] |= overflow;
                }
                ++j;
            }
        }
        line.resize((sub.width + 7) / 8); // drop whole trailing bytes
    }
    return sub;
}

void frame::modify(int start, int end, const data_type &lines)
{
    for (int i = 0; i < data.size(); ++i)
    {
        int src_pos = 0;
        for (auto dest_pos = start; dest_pos < end;)
        {
            // every iteration covers the range of bits between two byte-borders
            // in both source and destination line arrays
            auto src_byte_pos = src_pos % 8; // counting from MSB to LSB
            auto src_bits = std::min(end - dest_pos, 8 - src_byte_pos);
            auto src_byte = lines[i][src_pos / 8];

            auto dest_byte_pos = dest_pos % 8;
            auto dest_bits = 8 - dest_byte_pos;
            auto &dest_byte = data[i][dest_pos / 8];

            auto bits = std::min(src_bits, dest_bits);
            auto mask = (1 << bits) - 1;
            auto dest_shift = 8 - dest_byte_pos - bits;
            auto src_shift = 8 - src_byte_pos - bits;

            src_byte = (src_byte >> src_shift) & mask;
            dest_byte &= ~(mask << dest_shift);  // clear bits in dest byte
            dest_byte |= src_byte << dest_shift; // apply bits from src byte

            src_pos += bits;
            dest_pos += bits;
        }
    }
}

void frame::text(int pos, std::u16string str)
{
    for (auto c : str)
    {
        modify(pos, pos + font::width, font::get(c));
        pos += font::width;
    }
}
