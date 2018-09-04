#include "atb.h"

#include <boost/beast/core.hpp>
#include <boost/beast/http.hpp>
#include <boost/beast/version.hpp>
#include <boost/asio/connect.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/asio/ssl/error.hpp>
#include <boost/asio/ssl/stream.hpp>

#include <boost/property_tree/ptree.hpp>
#include <boost/property_tree/json_parser.hpp>

#include <boost/date_time.hpp>

#include <locale>
#include <string>
#include <sstream>

namespace http = boost::beast::http;

static const std::string host = "atbapi.tar.io";
static const std::string port = "443";
static const std::string api_path = "/api/v1/departures/";

struct departure
{
    std::string line;
    boost::posix_time::ptime scheduled_departure_time;
    std::string destination;
    bool is_realtime_data;
};

static http::dynamic_body request_departures(std::string node_id);
static std::vector<departure> parse_departures(std::istream &resp);
static boost::posix_time::ptime parse_timestamp(std::string s);

std::string get_bus_sign(std::string node_id, std::initializer_list<std::string> lines,
                         std::chrono::milliseconds pre_wait)
{
    auto response = request_departures(node_id);
    std::stringstream ss;
    ss << boost::beast::buffers(response);
    auto departures = parse_departures(ss);

    auto now = boost::posix_time::second_clock::local_time();
    std::string result;
    for (auto &line : lines)
    {
        if (!result.empty())
            result += " | ";
        result += line + "->";
        int found = 0;
        for (const auto &departure : departures)
        {
            if (departure.line == line)
            {
                auto d = departure.scheduled_departure_time - now;
                auto m = (d.total_seconds() + 30) / 60;
                if (found == 0)
                    result += departure.destination + std::to_string(m) + 'm';
                else
                    result += " (" + std::to_string(m) + "m)";
                if (++found >= 2)
                    break;
            }
        }
        if (found == 0)
            result += '?';
    }
    return result;
}

static http::dynamic_body request_departures(std::string node_id)
{
    auto const target = api_path + node_id;

    // The io_context is required for all I/O
    boost::asio::io_context ioc;

    // The SSL context is required, and holds certificates
    namespace ssl = boost::asio::ssl;
    ssl::context ctx{ssl::context::sslv23_client};

    // These objects perform our I/O
    using tcp = boost::asio::ip::tcp;
    tcp::resolver resolver{ioc};
    ssl::stream<tcp::socket> stream{ioc, ctx};

    // Set SNI Hostname (many hosts need this to handshake successfully)
    if (!SSL_set_tlsext_host_name(stream.native_handle(), host))
    {
        boost::system::error_code ec{static_cast<int>(::ERR_get_error()), boost::asio::error::get_ssl_category()};
        throw boost::system::system_error{ec};
    }

    // Look up the domain name
    auto const results = resolver.resolve(host, port);

    // Make the connection on the IP address we get from a lookup
    boost::asio::connect(stream.next_layer(), results.begin(), results.end());

    // Perform the SSL handshake
    stream.handshake(ssl::stream_base::client);

    // Set up an HTTP GET request message
    http::request<http::string_body> req{http::verb::get, target, 11};
    req.set(http::field::host, host);
    req.set(http::field::user_agent, BOOST_BEAST_VERSION_STRING);

    // Send the HTTP request to the remote host
    http::write(stream, req);

    // This buffer is used for reading and must be persisted
    boost::beast::flat_buffer buffer;

    // Declare a container to hold the response
    http::response<http::dynamic_body> res;

    // Receive the HTTP response
    http::read(stream, buffer, res);

    if (res.result_int() != 200)
        throw std::runtime_error("Error returned by server");

    // Gracefully close the stream
    boost::system::error_code ec;
    stream.shutdown(ec);

    return res.body();
}

static std::vector<departure> parse_departures(std::istream &resp)
{
    namespace pt = boost::property_tree;

    std::vector<departure> departures;
    pt::ptree root;
    pt::read_json(resp, root);

    for (const auto &d : root.get_child("departures"))
    {
        departures.push_back(departure{
            .line = d.second.get<std::string>("line"),
            .scheduled_departure_time = parse_timestamp(d.second.get<std::string>("scheduledDepartureTime")),
            .destination = d.second.get<std::string>("destination"),
            .is_realtime_data = d.second.get<bool>("isRealtimeData"),
        });
    }

    return departures;
}

static boost::posix_time::ptime parse_timestamp(std::string s)
{
    std::stringstream ss(s);
    ss.exceptions(std::ios_base::failbit);

    auto facet = new boost::posix_time::time_input_facet("%Y-%m-%dT%H:%M:%S%F%Q%F *");
    ss.imbue(std::locale(ss.getloc(), facet));

    boost::posix_time::ptime t;
    ss >> t;
    return t;
}
