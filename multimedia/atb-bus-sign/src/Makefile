LIBS = -lssl -lcrypto -lboost_system
SRCS = main.cpp atb.cpp font.cpp frame.cpp max7219.cpp
OBJS = $(SRCS:.cpp=.o)
SRCS_SIM = main_sim.cpp atb.cpp font.cpp frame.cpp
OBJS_SIM = $(SRCS_SIM:.cpp=.o)

atb-bus-sign: $(OBJS)
	$(CXX) -o $@ $^ $(LDFLAGS) $(LIBS) -pthread

sim: $(OBJS_SIM)
	$(CXX) -o $@ $^ $(LDFLAGS) $(LIBS) -pthread

%.o: %.cpp
	$(CXX) -c -o $@ $< $(CXXFLAGS) -pthread

.PHONY: clean
clean:
	rm -f $(OBJS) $(OBJS_SIM) atb-bus-sign sim
