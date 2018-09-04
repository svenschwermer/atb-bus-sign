#include "frame.h"
#include "font.h"

#include <algorithm>
#include <codecvt>
#include <iomanip>
#include <iterator>
#include <locale>
#include <sstream>

frame::frame(std::string str, size_t minimum_width)
    : data(font::height)
{
    auto utf16 = std::wstring_convert<std::codecvt_utf8_utf16<char16_t>, char16_t>{}.from_bytes(str.data());
    width = std::max(utf16.length() * font::width, minimum_width);
    for (auto &line : data)
        line.resize(width);
    text(0, utf16);
}

frame::frame(int lines, int columns)
    : data(lines)
{
    for (auto &line : data)
        line.resize(columns);
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
    auto sub = duplicate();
    auto w = end - start;
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
                if (j > 0)
                {
                    auto overflow = (b >> (8 - shift)) & mask;
                    line[j - 1] |= overflow;
                }
                b <<= shift;
                ++j;
            }
        }
        line.resize((w + 7) / 8); // drop whole trailing bytes
    }
    sub.width = w;
    return sub;
}

std::string frame::string() const
{
    std::stringstream ss;
    for (int i = 0; i < width; ++i)
        ss << std::setw(2) << i << ' ';
    ss << '\n';
    for (auto &line : data)
    {
        int j = 0;
        for (auto pixels : line)
        {
            for (int i = 0; i < 8 && j * 8 + i < width; ++i)
            {
                if (pixels & (1 << (7 - i)))
                    ss << u8" █ ";
                else
                    ss << " . ";
            }
            ++j;
        }
        ss << '\n';
    }
    return ss.str();
}

frame frame::duplicate() const
{
    auto dup = frame(data.size(), 2 * width);
    for (int i = 0; i < data.size(); ++i)
    {
        std::copy(data.begin(), data.end(), dup.data.begin());
        dup.modify(width, 2 * width, data);
    }
    return dup;
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

std::ostream &operator<<(std::ostream &os, const frame &f)
{
    return os << f.string();
}
