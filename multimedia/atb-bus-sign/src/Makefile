CFLAGS = -O2 -pthread
LIBS = -lssl -lcrypto -lboost_system
LDFLAGS = $(LIBS) $(CFLAGS) -s -static

SRCS = main.cpp atb.cpp font.cpp frame.cpp max7219.cpp
OBJS = $(SRCS:.cpp=.o)

atb-bus-sign: $(OBJS)
	$(CXX) -o $@ $^ $(LDFLAGS)

%.o: %.cpp
	$(CXX) -c -o $@ $< $(CFLAGS)
