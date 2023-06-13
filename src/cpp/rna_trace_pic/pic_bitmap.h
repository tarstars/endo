#ifndef __PIC_BITMAP_H__INCLUDED__
#define __PIC_BITMAP_H__INCLUDED__

#include <cstdint>
#include <vector>

struct CColor
{
    CColor(uint8_t a_,
           uint8_t r_,
           uint8_t g_,
           uint8_t b_ ) :
        b(b_), g(g_), r(r_), a(a_)
    {}

    bool operator ==(CColor const& c) const
    {
        return a==c.a &&
            b==c.b &&
            g==c.g &&
            r==c.r;
    }

    void compose( CColor p )
    {
        a = p.a + (int(a)*(255-p.a)/255);
        r = p.r + (int(r)*(255-p.a)/255);
        g = p.g + (int(g)*(255-p.a)/255);
        b = p.b + (int(b)*(255-p.a)/255);
    }

    void clip( CColor p )
    {
        a = int(a)*p.a/255;
        r = int(r)*p.a/255;
        g = int(g)*p.a/255;
        b = int(b)*p.a/255;
    }

    uint8_t b;
    uint8_t g;
    uint8_t r;
    uint8_t a;
} __attribute__((__packed__));

class CBitmap
{
public:
    enum Consts
    {
        WIDTH = 600,
        HEIGHT = 600,
    };
    const int LENGTH = WIDTH*HEIGHT;

    CBitmap(void);
    CColor* row(int ind);

    CColor dat_[WIDTH*HEIGHT];
};

CBitmap do_work(void);

#endif //__PIC_BITMAP_H__INCLUDED__
