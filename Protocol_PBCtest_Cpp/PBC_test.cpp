//高精度时控函数QueryPerformanceFrequency（），QueryPerformanceCounter（）
#include<windows.h>
#include<iostream>

#include <pbc.h>
#include <pbc_test.h>
#define LEN 6

using namespace std;
int main(int argc, char **argv)
{
    double time=0;//用来计时，精确到微秒
	double totaltime=0;//用来计时，精确到微秒
	double maxtime=0;//用来计时，精确到微秒
	double minitime=100;//用来计时，精确到微秒
    //定义双线性对运算的参数
    pairing_t pairing;
	element_t s,x,r; //整数
	element_t P,Ppub,Qu,Du,Su,Xu,Yu,V; //G1上的群元素
	element_t T1,T2;//GT上的群元素
	//double time1,time2;
	int byte;
	pbc_demo_pairing_init(pairing, argc, argv);
    
	//将变量初始化为Zr上的元素
	element_init_Zr(s,pairing);
	element_init_Zr(r,pairing);
	element_init_Zr(x,pairing);
    
	//将变量初始化为G1上的元素
	element_init_G1(P,pairing);
	element_init_G1(Ppub,pairing);
	element_init_G1(Xu,pairing);
	//element_init_G1(Yu,pairing);
	
	//将变量初始化为GT中的元素
	element_init_GT(T1,pairing);
	//element_init_GT(T2,pairing);

    
	//判断所用的配对是否为对称配对
	if(!pairing_is_symmetric(pairing)){
		fprintf(stderr,"只能在对称配对下运行");
			exit(1);
    }
//-------------------------------------------
    //element_random(s);//随机选择s
	//element_random(P);//随机选择P
    //element_mul_zn(Ppub,P,s); //公钥就是做了一次标量乘法Ppub=sP
    
    //element_random(x);//随机选择x
	//element_mul_zn(Xu,P,x);//Xu=xP
	//element_mul_zn(Yu,Ppub,x);//Yu=xP    
			//element_random(s);
			//element_random(P);
        //element_random(Ppub);
//-----------test-------computation_cost----------------
	int j=0;
    for(j;j<100;j++){
		element_random(Ppub);
		element_random(Xu);
	LARGE_INTEGER nFreq;
	LARGE_INTEGER nBeginTime;
	LARGE_INTEGER nEndTime;
	QueryPerformanceFrequency(&nFreq);
        QueryPerformanceCounter(&nBeginTime);//开始计时
    //for中是测试内容
	    //for(int i=0;i<1;i++)
        //{
			//写个循环测试多次取平均值，减少误差
			//element_random(s);
			//element_random(P);
			//element_mul_zn(Ppub,P,s); //公钥就是做了一次标量乘法Ppub=sP
			//element_random(Xu);
			//element_random(Ppub);
            pairing_apply(T1,Xu,Ppub,pairing);//T1=e(Xu,Ppub)
            //pairing_apply(T2,Yu,P,pairing);//T2=e(Yu,P)
        //}
        QueryPerformanceCounter(&nEndTime);//停止计时
        time=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为s
        cout<<"程序执行时间："<<time*1000<<"ms"<<endl;
		if(maxtime<time)
		{maxtime=time;}

		if(minitime>time)
		{minitime=time;}
		totaltime = totaltime + time;
		
	}   
	cout<<"程序执行总时间："<<totaltime*1000<<"ms"<<endl;
	cout<<"程序执行平均时间："<<totaltime /(j) *1000<<"ms"<<endl;
	cout<<"程序执行最大时间："<<maxtime *1000<<"ms"<<endl;
	cout<<"程序执行最小时间："<<minitime *1000<<"ms"<<endl;
    return 0;
}
