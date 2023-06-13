//#include "SDL.h"    
#include <CImg.h>
#include <stdio.h> 
#include <cmath>
#include <algorithm>
#include <fstream>
#include <iostream>

#include "picBitmap.h"
#include "canvas_out.h"

using namespace std;
using namespace cimg_library;

typedef CImg<unsigned char> CBmp;

class CScene
{
	CCanvasOut sam;
	int stepNum;
	bool isGoForward;
public:
	CScene(void);
	void DrawPicture(CBmp *,int ,int );
	void tryAddBmp(void);
	void goForward(void){isGoForward^=true;}
};


CScene::CScene(void):stepNum(0),isGoForward(false)
{
}

void 
CScene::DrawPicture(CBmp *screen,int mx,int my)
{
    static bool theend = true;

    for(int t=0;t<100 && sam.canContinue() && (!sam.isSkipAddBmp() || isGoForward);++t)
    {
        sam.step();
        stepNum++;
        if(isGoForward&&sam.isSkipAddBmp())
        {
            //cout<<stepNum<<endl;
            sam.addBmp();
        }
    }

    if(!sam.canContinue() && theend)
    {
        cout << "THE END." << endl;
        theend = false;
    }
    int x,y;

    CBitmap bmp=sam.getPicture();
    for(x=0;x<600;++x)
        for(y=0;y<600;++y)
        {
            CColor pix=bmp[y][x];
			unsigned char color[]={pix.r,pix.g,pix.b};
			screen->draw_point(x,y,color);
        }
}

void
CScene::tryAddBmp(void)
{
	if (sam.isSkipAddBmp())
        {
            cout<<stepNum<<endl;
			sam.addBmp();
        }
}

int main() 
{
 	int mx=600,my=600;
	CBmp bmp(mx,my,1,3);
	CImgDisplay disp(bmp,"Endo forever!");
	bool keypressed=false;
	CScene myScene;

    while(!disp.is_closed()&&!keypressed)
	{
		switch(disp.key())
		{
			case cimg::keyESC: keypressed=true; break;
			case cimg::keySPACE: myScene.tryAddBmp(); break;
			case cimg::keyF: myScene.goForward(); break;
		}
		myScene.DrawPicture(&bmp,mx,my);
		disp.display(bmp);
	}
	bmp.save_png("res.png");
     
    return 0;
}

