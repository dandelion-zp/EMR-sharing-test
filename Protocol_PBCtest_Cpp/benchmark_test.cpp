#include<windows.h>
#include<iostream>

#include <stdint.h> // for intptr_t
#include "pbc.h"
#include "pbc_test.h"

#define LEN 6
/* I've heard that sometimes automatic garbage collection can outperform
 * manual collection, so I briefly tried using the Boehm-Demers-Weiser GC
 * library. Both GMP and PBC support custom memory allocation routines so
 * incorporating the GC library is trivial.
 *
 * Automatic garbage collection appears to slow this program down a little,
 * even if only PBC collects automatically. (The case where PBC collects
 * manually but GMP collects automatically cannot be achieved with the GC
 * library because PBC objects point at GMP objects.)
 *
 * Perhaps specially-tailored memory allocation routines could shave off
 * some time, but one would have to thoroughly analyze PBC and GMP memory usage
 * patterns.
 *
 * Below is the commented-out code that collects garbage for PBC. Of course,
 * if you want to use it you must also tell the build system where to find
 * gc.h and to link with the GC library.
 *
 * Also, you may wish to write similar code for GMP (which I unfortunately
 * deleted before thinking that it might be useful for others).
 * Note GC_MALLOC_ATOMIC may be used for GMP since the mpz_t type does not
 * store pointers in the memory it allocates.
 *
 * The malloc and realloc functions should exit on failure but I didn't
 * bother since I was only seeing if GC could speed up this program.

#include <gc.h>
#include <pbc_utils.h>

void *gc_alloc(size_t size) {
  return GC_MALLOC(size);
}

void *gc_realloc(void *ptr, size_t size) {
  return GC_REALLOC(ptr, size);
}

void gc_free(void *ptr) {
  UNUSED_VAR(ptr);
}

 * The following should be the first two statements in main()

GC_INIT();
pbc_set_memory_functions(gc_alloc, gc_realloc, gc_free);

 */

int main(int argc, char **argv) {
  pairing_t pairing;
  element_t x, y, r, r2, s, h;
  element_t P,Ppub,P1,P2; //G1上的群元素
  element_t T1,T2;//GT上的群元素
  int i, n;
  double t0, t1, ttotal, ttotalpp, ttotalSM, ttotalME, ttotalPA, ttotalHF;
  pairing_pp_t pp;

  // Cheat for slightly faster times:
  // pbc_set_memory_functions(malloc, realloc, free);

  pbc_demo_pairing_init(pairing, argc, argv);

  element_init_G1(x, pairing);
  element_init_G2(y, pairing);
  element_init_GT(r, pairing);
  element_init_GT(r2, pairing);

  element_init_Zr(h, pairing);

  element_init_Zr(s,pairing);
  element_init_G1(P, pairing);
  element_init_G1(P1, pairing);
  element_init_G1(P2, pairing);
  element_init_G1(Ppub, pairing);
  element_init_GT(T1, pairing);
  element_init_GT(T2, pairing);

  LARGE_INTEGER nFreq;
	LARGE_INTEGER nBeginTime;
	LARGE_INTEGER nEndTime;
	QueryPerformanceFrequency(&nFreq);

  n = 100;
  ttotal = 0.0;
  ttotalpp = 0.0;
  ttotalSM = 0.0;
  ttotalME = 0.0;
  ttotalPA = 0.0;
  ttotalHF = 0.0;
  for (i=0; i<n; i++) {
    element_random(x);
    element_random(y);

	element_random(s);
	element_random(P1);
	element_random(P2);
	element_random(T2);
	//------------Bilinear pairing-------------
    pairing_pp_init(pp, x, pairing); //x是G1群的元素，预处理
	QueryPerformanceCounter(&nBeginTime);//开始计时
    //t0 = pbc_get_time();
    pairing_pp_apply(r, y, pp); //带有预处理的对运算 r=e(x,y)
    QueryPerformanceCounter(&nEndTime);//停止计时
	//t1 = pbc_get_time();
	ttotalpp+=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms
    //ttotalpp += t1 - t0;
    pairing_pp_clear(pp); // don’t need pp anymore

    QueryPerformanceCounter(&nBeginTime);//开始计时
	//t0 = pbc_get_time();

    //element_pairing(r2, x, y); //普通双线性对运算 r2=e(x,y)
	pairing_apply(r2, x, y,pairing);//普通双线对运算 r2=e(x,y)
    QueryPerformanceCounter(&nEndTime);//停止计时
	//t1 = pbc_get_time();
	ttotal+=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms
    //ttotal += t1 - t0;

	//------------Scalar multiplication-------------
	QueryPerformanceCounter(&nBeginTime);//开始计时
	element_mul_zn(Ppub,P,s); //标量乘法 Ppub=sP
	QueryPerformanceCounter(&nEndTime);//停止计时
	ttotalSM +=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms

	//------------Modular exponentiation-------------
	QueryPerformanceCounter(&nBeginTime);//开始计时
	element_pow_zn(T1,T2,s); // 模幂运算 T1=T2^s
	QueryPerformanceCounter(&nEndTime);//停止计时
	ttotalME +=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms

	//------------Point addition-------------
	QueryPerformanceCounter(&nBeginTime);//开始计时
	element_add(P,P1,P2); //点加法 P=P1+P2
	QueryPerformanceCounter(&nEndTime);//停止计时
	ttotalPA +=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms

	//------------Hash function-------------
	QueryPerformanceCounter(&nBeginTime);//开始计时
	element_from_hash(h, "ABCDEF", 6); //hash映射
	//element_printf("h = %B\n", h);
	QueryPerformanceCounter(&nEndTime);//停止计时
	ttotalHF +=(double)(nEndTime.QuadPart-nBeginTime.QuadPart)/(double)nFreq.QuadPart;//计算程序执行时间单位为微秒。*1000转换成ms
	//------------------------------------------------------------------

    element_printf("x = %B\n", x);
    element_printf("y = %B\n", y);
    element_printf("e(x,y) = %B\n", r);
    if (element_cmp(r, r2)) {
      printf("BUG!\n");
      exit(1);
    }
  }
  printf("average bilinear pairing time = %f ms\n", ttotal*1000 / n);
  printf("average bilinear pairing time (preprocessed) = %f ms\n", ttotalpp*1000 / n);
  printf("average scalar multiplication time = %f ms\n", ttotalSM*1000 / n);
  printf("average modular exponentiation time = %f ms\n", ttotalME*1000 / n);
  printf("average point addition time = %f ms\n", ttotalPA*1000 / n);
  printf("average hash function time = %f ms\n", ttotalHF*1000 / n);

  element_clear(x);
  element_clear(y);
  element_clear(r);
  element_clear(r2);

  pairing_clear(pairing);

  return 0;
}
