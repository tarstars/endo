#include <iomanip>
#include <iostream>
#include <sstream>
#include <fstream>
#include <vector>
#include <stack>
#include <ext/rope>

using namespace std;
using namespace __gnu_cxx;

#define FORE(itx,a) for(typeof((a).begin()) itx = (a).begin(); itx!=(a).end(); ++itx)
#define FORR(x,a,b) for(int x=a;x<b;++x)

// #define SHOW_EVERY 100000
// #define DUMP_EVERY 100000

char pop(__gnu_cxx::crope &dna,int &ind)
{
    if (ind>=dna.size()) throw(std::string("dna empty"));
	return dna[ind++];
}

int nat(crope &dna, int& ind)
{
	switch(pop(dna,ind))
	{
		case 'P': return 0;
		case 'I': case 'F': return 2*nat(dna,ind);
		case 'C': return 2*nat(dna,ind)+1;
	}
}

crope consts(crope& dna,int &ind)
{
	crope ret;	
	while(true){
		bool thatsAll=true;
		if (ind<dna.size())
			switch(dna[ind]){
				case 'C':++ind;ret+="I";thatsAll=false;break;
				case 'F':++ind;ret+="C";thatsAll=false;break;
				case 'P':++ind;ret+="F";thatsAll=false;break;
				}
		if(thatsAll && ind+1<dna.size() && dna[ind]=='I' && dna[ind+1]=='C')
		{
			ind+=2;
			thatsAll=false;
			ret+="P";
		}

		if(thatsAll)
			return ret;
	}
}

void make_pattern(crope& dna,int& ind,ostream &rna,ostream& dest)
{
	int lvl=0;
	char c;
	while(true)
	{
		switch(pop(dna,ind)){
			case 'C': dest<<"b I  ";break;
			case 'F': dest<<"b C  ";break;
			case 'P': dest<<"b F  ";break;
			case 'I':
				switch(pop(dna,ind)){
					case 'C':dest<<"b P  ";break;
					case 'P':dest<<"! "<<nat(dna,ind)<<"  ";break;
					case 'F':
							pop(dna,ind);
							dest<<"? "<<consts(dna,ind)<<"  ";
							break;
					case 'I':
						switch(pop(dna,ind)){
							case 'P':dest<<"(";++lvl;break;
							case 'C':
							case 'F':if(lvl==0) return; else --lvl;dest<<")";break;
							case 'I':
								FORR(t,0,7)
									rna<<pop(dna,ind);
                                                                rna<<endl;
								break;
							}
					}
			}
	}
}

void make_template(crope &dna,int &ind,ostream &rna,ostream& dest)
{
	while(true)
	{
		switch(pop(dna,ind)){
			case 'C': dest<<"b I  ";break;
			case 'F': dest<<"b C  ";break;
			case 'P': dest<<"b F  ";break;
			case 'I':
				switch(pop(dna,ind)){
					case 'C':dest<<"b P  ";break;
					case 'P':
					case 'F':{int l,n;l=nat(dna,ind);n=nat(dna,ind);dest<<"r "<<n<<" "<<l<<"  ";break;}
					case 'I':
						switch(pop(dna,ind)){
							case 'C':
							case 'F': return;
							case 'P': dest<<"m "<<nat(dna,ind)<<"  ";break;
							case 'I': 
								FORR(t,0,7) 
									rna<<pop(dna,ind); 
                                                                rna<<endl;
								break;
							}
					}
			}
	}
}

crope quote(crope d)
{
	crope ret;
	FORE(itx,d)
		switch (*itx)
		{
			case 'I': ret+="C";break;
			case 'C': ret+="F";break;
			case 'F': ret+="P";break;
			case 'P': ret+="IC";break;
		}
	return ret;
}

crope prot(int l, crope d)
{
	FORR(t,0,l)
		d=quote(d);
	return d;
}

crope asnat(int a)
{
	if(a==0) return "P";
	if(a%2==0) return crope("I")+asnat(a/2);
	if(a%2==1) return crope("C")+asnat(a/2);
    return "";
}

void match_replace(crope& dna,istream& sour_pat,istream& sour_tmpl,int stage)
{
	stack<int> stc;
	string curs;
	vector<crope> env;
	int i=0;
	char c;
	bool any=false;

    cerr << "i = " << i << endl;

	while(sour_pat>>c)
		switch(c)
		{
			case 'b': 
				sour_pat>>c; 
				if(i<int(dna.size()) && c==dna[i]) 
					++i; 
                else {
                    cerr << "match failed" << endl;
                    return;
                }
				break;
			case '!': {
				int n;sour_pat>>n; 
				i+=n; 
                if (i>int(dna.size())) {
                    std::cerr << "match failed" << endl;
					return; 
                }
                cerr << "skip to " << i << endl;
				break;}
			case '?':
				{
					string word;
					sour_pat>>word;
					size_t x=dna.find(word.c_str(),i);
                    if(x==string::npos) {
                        std::cerr << "match failed" << endl;
                        return;
                    }
                    size_t old_i = i;
					i=x+word.size();
                    std::cerr << "?: successfull advance to " << i << " from " << old_i << endl;
					break;
				}
			case '(': stc.push(i);break;
			case ')': env.push_back(dna.substr(stc.top(),i-stc.top()));stc.pop();break;
		}

	if (i<=int(dna.size())) {
	    cerr << "truncate dna from " << i << " position"  << endl;
		dna.erase(0,i);
	}

	crope r;

    cerr << "env size = " << env.size() << endl;
    for (int p = 0; p < env.size(); p++) {
        cerr << "\tenv[i].size = " << env[p].size() << endl;
    }


	while(sour_tmpl>>c)
		switch(c){
			case 'b':  sour_tmpl>>c; r+=c;break;
			case 'r': {int n,l;sour_tmpl>>n>>l;  r+=prot(l,(n<env.size())?env[n]:"");break;}
			case 'm':{int n;sour_tmpl>>n; r+=asnat((n<env.size())?env[n].size():0);break;}
		}
	dna=r+dna;
}

void do_work(const string& source, const string& destination)
{
	crope dna;
	{
        ifstream source_stream(source);
		string once_upon;
        source_stream>>once_upon;
		dna=crope(once_upon.c_str());
	}
    int ind=0;

	try{
        stringstream pat_mem,tem_mem;
        make_pattern(dna,ind,cout,pat_mem);
        make_template(dna,ind,cout,tem_mem);

        cerr<<"pat="<<pat_mem.str()<<endl;
        cerr<<"tem="<<tem_mem.str()<<endl;

        dna.erase(0,ind);

        match_replace(dna,pat_mem,tem_mem, 0);

        ofstream dump(destination.c_str());
        dump.write(dna.c_str(), dna.size());
        dump.close();
    }catch(string msg){cerr<<"End of program:"<<msg<<endl;}
}

void numberTest(void)
{
	bool isOk=true;
	int fn=0;
	for(int t=0;t<1000000&&isOk;++t)
	{
		crope sr=asnat(t);
		if(t%10000==0)
			cerr<<t<<" "<<sr<<endl;
		int ind=0,tt=nat(sr,ind);
		if (tt!=t)
			isOk=false,fn=t;
	}
	if (isOk)
	{
		cerr<<"number test passed"<<endl;
	}
	else
		cerr<<"number test failed in "<<fn<<" case"<<endl;
}

void stringTest(void)
{
	bool isOk=true;
	int fn=0;
	crope joke("IPF");

	for(int t=100000;t<2000000&&isOk;t+=100000)
	{
		crope st=asnat(t);
		crope suff=crope("I")+joke[t%3]+crope("IF");
		crope pst=quote(st)+suff;
		int ind=0;
		crope sst=consts(pst,ind);
		if (t%100000==0)
			cerr<<"st="<<st<<"\tpst="<<pst<<"\tsst="<<sst<<endl;
		if (sst!=st || pst.substr(ind,pst.size()-ind)!=suff)
			isOk=false,fn=t;			
	}

	if (isOk)
		cerr<<"string test passed"<<endl;
	else
		cerr<<"string test failed in "<<fn<<" case"<<endl;
}

void stringTest00(void)
{
	int ind;
	crope joke("IPF");
	for(int t=0;t<3;++t)
	{
		ind=0;
		crope st=asnat(100000);
		crope dna=quote(st)+crope("I")+joke[t]+crope("F");
		cerr<<"dna="<<dna<<endl;
		crope sst=consts(dna,ind);
		cerr<<" st="<<st<<endl;
		cerr<<"sst="<<sst<<endl;
		cerr<<"remind="<<dna.substr(ind,dna.size()-ind)<<endl<<endl;
	}
}

void doTests(void)
{
	//numberTest();
	stringTest();
	//stringTest00();
}

int main(int argc, char ** argv)
{
    do_work(string(argv[1]), string(argv[2]));
	//doTests();
	return 0;
}

