#include <string>
#include <iostream>
#include <fstream>
#include <map>
#include <vector>
#include <deque>
#include <stack>

#include "picBitmap.h"
#include "canvas_out.h"

using namespace std;

CColor::CColor(int xr,int xg,int xb,int xa):r(xr),g(xg),b(xb),a(xa){}

bool 
CColor::operator==(const CColor& nulPlane)const
{
	return r==nulPlane.r && g==nulPlane.g && b==nulPlane.b && a==nulPlane.a;
}

void 
CColor::compose(const CColor& nulPlane)
{
(*this)=CColor(
		nulPlane.r+r*(255-nulPlane.a)/255,
		nulPlane.g+g*(255-nulPlane.a)/255,
		nulPlane.b+b*(255-nulPlane.a)/255,
		nulPlane.a+a*(255-nulPlane.a)/255
		);
}

void CColor::clip(const CColor& nulPlane)
{
(*this)=CColor(	
		r*nulPlane.a/255,
		g*nulPlane.a/255,
		b*nulPlane.a/255,
		a*nulPlane.a/255);
}

class CBuck
{
public:
	int sr,sg,sb,sa,nc,na;
	CBuck():sr(0),sg(0),sb(0),sa(0),nc(0),na(0){}
	void clear(){*this=CBuck();}
	void aC(int r,int g,int b){sr+=r;sg+=g;sb+=b;++nc;}
	void aA(int a){sa+=a;++na;}
	CColor curPix(void)
	{
		int ac=255;
		if(na)
			ac=sa/na;

		if (nc) 
			return CColor(sr/nc*ac/255,sg/nc*ac/255,sb/nc*ac/255,ac); 
		else 
			return CColor(0,0,0,ac);
	}
};


CBitmap::CBitmap(void):dat(600*600){}
VCCI CBitmap::operator[](int ind){return dat.begin()+600*ind;}

typedef deque<CBitmap> DCB;

class CCanvas
{
	DCB dat;
	CBuck buck;
	int posx,posy,markx,marky,dir;
	static int dirdx[],dirdy[];
	map<string,int> opcode;
	int t;

public:
	CCanvas();
	void line(int x1,int y1,int x2,int y2);
	void tryFill(void);
	void addBmp(void){if(dat.size()<10) dat.push_front(CBitmap());}
	void compose(void);
	void clip(void);

	void setPixel(int lx,int ly){dat[0][ly][lx]=buck.curPix();}
	CColor getPixel(int lx,int ly){return dat[0][ly][lx];}
	void fill(int lx,int ly,CColor& init);
	void step(void);
	CBitmap getPicture(void);
	bool canContinue(void)const{return cin.good();}

	bool skipAddBmp;
	void doAddBmp(void);

    int rnaMeter;
};

int CCanvas::dirdx[]={1,0,-1,0};
int CCanvas::dirdy[]={0,-1,0,1};

CCanvas::CCanvas():posx(0),posy(0),markx(0),marky(0),dir(0),t(0),skipAddBmp(false),rnaMeter(0)
{
	dat.push_front(CBitmap());

	opcode["PIPIIIC"]=1;
	opcode["PIPIIIP"]=2;
	opcode["PIPIICC"]=3;
	opcode["PIPIICF"]=4;
	opcode["PIPIICP"]=5;
	opcode["PIPIIFC"]=6;
	opcode["PIPIIFF"]=7;
	opcode["PIPIIPC"]=8;
	opcode["PIPIIPF"]=9;
	opcode["PIPIIPP"]=10;
	opcode["PIIPICP"]=11;
	opcode["PIIIIIP"]=12;
	opcode["PCCCCCP"]=13;
	opcode["PFFFFFP"]=14;
	opcode["PCCIFFP"]=15;
	opcode["PFFICCP"]=16;
	opcode["PIIPIIP"]=17;
	opcode["PCCPFFP"]=18;
	opcode["PFFPCCP"]=19;
	opcode["PFFICCF"]=20;
}

int fabs(int x){return (x<0)?-x:x;}

void
CCanvas::line(int x0,int y0,int x1,int y1)
{
	int dx(x1-x0),dy(y1-y0);
	int d(max(fabs(dx),fabs(dy)));
	int c=0;
	if(dx*dy<=0) c=1;
	int x(d*x0+(d-c)/2),y(y0*d+(d-c)/2);
	for(int t=0;t<d;++t)
	{
		setPixel(x/d,y/d);
		x+=dx;
		y+=dy;
	}
	setPixel(x1,y1);
}

void
CCanvas::tryFill(void)
{
	CColor nw=buck.curPix(),old=getPixel(posx,posy);
	if(!(nw==old))
		fill(posx,posy,old);
}

class CPoint
{
public:
	int x,y;
	CPoint(int xx=0,int xy=0):x(xx),y(xy){}	
};

typedef stack<CPoint> SCP;

void
CCanvas::fill(int lx,int ly,CColor& init)
{
	SCP stor;
	stor.push(CPoint(lx,ly));
	while(!stor.empty())
	{		
		CPoint cur=stor.top();stor.pop();
		lx=cur.x;ly=cur.y;
		setPixel(lx,ly);
		if(lx>0 && getPixel(lx-1,ly)==init) stor.push(CPoint(lx-1,ly));
		if(lx<599 && getPixel(lx+1,ly)==init) stor.push(CPoint(lx+1,ly));
		if(ly>0 && getPixel(lx,ly-1)==init) stor.push(CPoint(lx,ly-1));
		if(ly<599 && getPixel(lx,ly+1)==init) stor.push(CPoint(lx,ly+1));
	}
}

void
CCanvas::compose(void)
{
    cerr << "compose at " << rnaMeter << " step" << endl;
	if (dat.size()<2) return;
	for(int y=0;y<600;++y)
		for(int x=0;x<600;++x)
			(dat[1])[y][x].compose((dat[0])[y][x]);
	dat.pop_front();
}

void
CCanvas::clip(void)
{
	if (dat.size()<2) return;
	for(int y=0;y<600;++y)
		for(int x=0;x<600;++x)
			(dat[1])[y][x].clip((dat[0])[y][x]);
	dat.pop_front();
}

void CCanvas::step(void)
{
    string word;
    if(!(cin>>word)) return;
	switch(opcode[word])
	{
		case 1: buck.aC(0,0,0);break;
		case 2: buck.aC(255,0,0);break;
		case 3: buck.aC(0,255,0);break;
		case 4: buck.aC(255,255,0);break;
		case 5: buck.aC(0,0,255);break;
		case 6: buck.aC(255,0,255);break;
		case 7: buck.aC(0,255,255);break;
		case 8: buck.aC(255,255,255);break;
		case 9 : buck.aA(0);break;
		case 10: buck.aA(255);break;
		case 11: buck.clear();break;
		case 12: posx=(posx+600+dirdx[dir])%600;posy=(posy+600+dirdy[dir])%600;break;
		case 13: dir=(dir+1)%4;break;
		case 14: dir=(dir+3)%4;break;
		case 15: markx=posx;marky=posy;break;
		case 16: line(posx,posy,markx,marky);break;
		case 17: tryFill();break;
		case 18: skipAddBmp=true;break;
		case 19: compose();break;
		case 20: clip();break;
		default: cerr<<"unrecognized rna: "<<word<<endl;
	}	
    ++rnaMeter;
}

void
CCanvas::doAddBmp(void){addBmp();skipAddBmp=false;}

CBitmap CCanvas::getPicture(void)
{
	return dat[0];
}


CCanvasOut::CCanvasOut(void)
{
	pDat=new CCanvas;
}

CCanvasOut::~CCanvasOut()
{
	delete pDat;
}

bool
CCanvasOut::canContinue(void)
{
	return pDat->canContinue();
}

void
CCanvasOut::step(void)
{
	pDat->step();
}

CBitmap
CCanvasOut::getPicture(void)
{
	return pDat->getPicture();
}

bool
CCanvasOut::isSkipAddBmp(void)
{
	return pDat->skipAddBmp;
}

void
CCanvasOut::addBmp(void)
{
	pDat->doAddBmp();
}

