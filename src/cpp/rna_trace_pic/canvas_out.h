#pragma once

class CBitmap;
class CCanvas;

class CCanvasOut
{
	CCanvas *pDat;
public:
	CCanvasOut(void);
	~CCanvasOut();
	bool canContinue(void);
	void step(void);
	CBitmap getPicture(void);
	bool isSkipAddBmp(void);
	void addBmp(void);
};


