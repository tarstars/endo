CXXFLAGS=-g -O3
CPPFLAGS=-DSHOW_EVERY=1
CPPFLAGS_DUMP=$(CPPFLAGS) -DDUMP_EVERY=1
SRCS=rna_gen.cpp
OBJS=$(SRCS:%.cpp=%.o)
OBJS_DUMP=$(SRCS:%.cpp=%_dump.o)

all: rna_gen rna_gen_dump

install: rna_gen rna_gen_dump
	cp $^ ../../../bin

clean:
	rm -f rna_gen rna_gen_dump $(OBJS) $(OBJS_DUMP)

rna_gen: $(OBJS)
	$(CXX) $(LDFLAGS) -o $@ $^

rna_gen_dump: $(OBJS_DUMP)
	$(CXX) $(LDFLAGS) -o $@ $^

%_dump.o: %.cpp
	$(CXX) -c $(CPPFLAGS_DUMP) $(CXXFLAGS) -o $@ $<
