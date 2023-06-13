#pragma once

#include <vector>

using namespace std;

class CColor
{
public:
    unsigned char r,g,b,a;
	CColor(int xr=0,int xg=0,int xb=0,int xa=0);
	bool operator==(const CColor& ri)const;
	void compose(const CColor& ri);
	void clip(const CColor& ri);
};

typedef vector<CColor> VCC;
typedef VCC::iterator VCCI;

class CBitmap
{
public:
	VCC dat;
	CBitmap(void);
	VCCI operator[](int ind);
};

CBitmap do_work(void);
