LIBS = -lssl -lcrypto -lboost_system
SRCS = main.cpp atb.cpp font.cpp frame.cpp max7219.cpp
OBJS = $(SRCS:.cpp=.o)

atb-bus-sign: $(OBJS)
	$(CXX) -o $@ $^ $(LDFLAGS) $(LIBS) -pthread

%.o: %.cpp
	$(CXX) -c -o $@ $< $(CXXFLAGS) -pthread
