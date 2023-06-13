TARGET = show_pic
LIBS =-L/usr/X11R6/lib -lX11 -lpthread
 
CFLAGS = -g -lm  -O3
CC = g++


all: $(TARGET)

$(TARGET): cimg_part.cpp mkpic.cpp
	$(CC) $(CFLAGS) -c cimg_part.cpp  -o cimg_part.o
	$(CC) $(CFLAGS) -c mkpic.cpp -o mkpic.o
	$(CC) -g -o $(TARGET) cimg_part.o mkpic.o $(LIBS)
install: $(TARGET)
	cp $(TARGET) ../../../bin
clean:
	rm *.o 
	rm $(TARGET)
