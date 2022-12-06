package main

import (
	"fmt"
	"github.com/Nik-U/pbc"
	"time"
	"os"
	"encoding/pem"
	"crypto/x509"
	"crypto/sha512"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
)


var P, g, p, q, h, r2, u, s, P_pub, Qu, Su, Qu_1, Su_1, Q_HN, S_HN, Q_HN_1, S_HN_1, R *pbc.Element

var Qc, Sc, Qc_1, Sc_1, QA, SA, QA_1, SA_1, Kca, Tc1, Tc2, Ta1, Ta2, Kac, sk_AC, sk_CA, sa, sb, a1, a2, sa_j *pbc.Element

var Qb, Sb, Qb_1, Sb_1, Tb1, Tb2, Kbc, Kcb, sk_BC, sk_CB, Ac, CP, LB *pbc.Element

var ID_u, ID_HN, ID_DN, Date, ID_c, ID_A, ID_b, m string

var j int = 10

var pairing *pbc.Pairing
//var x, X, aa [1000]*pbc.Element
var tt,rt [11]time.Duration

var rtmin, rtmax [11]time.Duration

func main() {

	//设置执行时间的最大值，最小值边界
	for i := 0; i < len(rt); i++ {
		rtmax[i] = 1 * time.Nanosecond
		rtmin[i] = 10000 * time.Millisecond
	}

	//设置程序执行的次数
	rp := 10000

	var j int

	for i := 0; i < rp; i++ {

		//1.系统初始化阶段
		SystemInit()

		//2.用户注册阶段
		Registration()

		//3.医疗数据管理阶段
		Med_manage()

		//4.医疗数据共享阶段
		Med_sharing()

		for j = 0; j< 11; j++ {
			if tt[j] > rtmax[j] {
				rtmax[j] = tt[j]
			}
			if tt[j] < rtmin[j] {
				rtmin[j] = tt[j]
			}

			rt[j] = rt[j] + tt[j]
			//fmt.Printf("执行 %v 次 rt[%v] = %v \n", rp, j, rt[j])
		}
	}

	var Ave, Avems [11]float64
	for k := 0; k < 11; k++ {
		Ave[k] = rt[k].Seconds() / float64(rp)
		Avems[k] = Ave[k] * 1000

		//seconds := 10
		//fmt.Print(time.Duration(seconds)*time.Second) // prints 10s

		fmt.Printf("执行 %v 次 rt[%v] = %v \n", rp, k, rt[k])
		fmt.Printf("rtmax = %v \n", rtmax[k])
		fmt.Printf("rtmin = %v \n", rtmin[k])
		fmt.Printf("rtave = %.6vms \n", Avems[k])
	}

}

//1.系统初始化，生成需要的参数
func SystemInit() {
	fmt.Println("----------System initialization 阶段----------")
	begin := time.Now()

	// In a real application, generate this once and publish it
	// GenerateA generates a pairing on the curve y^2 = x^3 + x over the field F_q
	// for some prime q = 3 mod 4. Type A pairings are symmetric (i.e., G1 == G2).
	// Type A pairings are best used when speed of computation is the primary
	// concern.
	// The authority generates system parameters
	params := pbc.GenerateA(256, 256)
	pairing = params.NewPairing()
	S := pairing.IsSymmetric()
	fmt.Printf("the pairing is symmetric: %t \n", S)

	// Initialize group elements. pbc automatically handles garbage collection.
	q = pairing.NewZr().Rand()

	P = pairing.NewG1().Rand()
	g = pairing.NewGT().Pair(P, P)

	//私钥参数
	s = pairing.NewZr().Rand()

	//公钥参数
	P_pub = pairing.NewG1().MulZn(P, s)


	tt[0] = time.Since(begin)
	fmt.Printf("the time of System initialization is tt[0]: %v \n", tt[0])

	fmt.Printf("P = %s\n", P)
	fmt.Printf("g = %s\n", g)

}

func Registration() {
	fmt.Println("----------Registration 阶段----------")
	begin := time.Now()

	//用户U注册
	ID_u = "1010010436"
	Qu = pairing.NewZr().SetFromHash([]byte(ID_u))

	Su0 := pairing.NewZr().Add(s, Qu)
	Su1 := pairing.NewZr().Invert(Su0)
	Su = pairing.NewG1().MulZn(P, Su1)

	Qu_1 = pairing.NewG1().SetFromHash([]byte(ID_u))
	Su_1 = pairing.NewG1().MulZn(Qu_1, s)

	//用户HN注册
	ID_HN = "0101020201"
	Q_HN = pairing.NewZr().SetFromHash([]byte(ID_HN))

	SHN0 := pairing.NewZr().Add(s, Q_HN)
	SHN1 := pairing.NewZr().Invert(SHN0)
	S_HN = pairing.NewG1().MulZn(P, SHN1)

	Q_HN_1 = pairing.NewG1().SetFromHash([]byte(ID_HN))
	S_HN_1 = pairing.NewG1().MulZn(Q_HN_1, s)

	r := pairing.NewZr().Rand()
	R0 := pairing.NewG1().MulZn(P, Q_HN)
	R1 := pairing.NewG1().Add(P_pub, R0)
	R = pairing.NewG1().MulZn(R1, r)

	ID_DN = "0202303002"
	Date = "2021-10-20"
	src1 := ID_HN + ID_DN + Date
	src := []byte(src1)

	sigText := SignatureRSA(src, "private.pem")

	tt[1] = time.Since(begin)
	fmt.Printf("the time of Registration is tt[1]: %v \n", tt[1])

	fmt.Printf("sigText = %s\n", sigText)

}

func Med_manage() {
	fmt.Println("----------Medical records management 阶段----------")
	begin := time.Now()

	//patient C 注册身份
	ID_c = "1010010413"
	Qc = pairing.NewZr().SetFromHash([]byte(ID_c))

	Sc0 := pairing.NewZr().Add(s, Qc)
	Sc1 := pairing.NewZr().Invert(Sc0)
	Sc = pairing.NewG1().MulZn(P, Sc1)

	Qc_1 = pairing.NewG1().SetFromHash([]byte(ID_c))
	Sc_1 = pairing.NewG1().MulZn(Qc_1, s)

	//hospital A 注册身份
	ID_A = "010102020A"
	QA = pairing.NewZr().SetFromHash([]byte(ID_A))

	SA0 := pairing.NewZr().Add(s, QA)
	SA1 := pairing.NewZr().Invert(SA0)
	SA = pairing.NewG1().MulZn(P, SA1)

	QA_1 = pairing.NewG1().SetFromHash([]byte(ID_A))
	SA_1 = pairing.NewG1().MulZn(QA_1, s)

	//patient C 生成的参数
	Kca = pairing.NewGT().Pair(Sc_1, QA_1)

	c1 := pairing.NewZr().Rand()
	c2 := pairing.NewZr().Rand()

	Tc1 = pairing.NewG1().MulZn(P, c1)
	Tc2 = pairing.NewG1().MulZn(P, c2)

	//hospital A 生成的参数
	a1 = pairing.NewZr().Rand()
	a2 = pairing.NewZr().Rand()

	Ta1 = pairing.NewG1().MulZn(P, a1)
	Ta2 = pairing.NewG1().MulZn(P, a2)

	Kac = pairing.NewGT().Pair(SA_1, Qc_1)

	fmt.Printf("Kac ?== Kca: %v \n", Kac.Equals(Kca))

	//A和C 协商秘钥
	aT1 := pairing.NewG1().MulZn(Tc1, a1).String()
	aT2 := pairing.NewG1().MulZn(Tc2, a2).String()
	Kac_str := Kac.String()
	Tc1_str := Tc1.String()
	Tc2_str := Tc2.String()
	Ta1_str := Ta1.String()
	Ta2_str := Ta2.String()

	str1 := ID_c + ID_A + aT1 + aT2 + Kac_str + Tc1_str + Tc2_str + Ta1_str + Ta2_str
	sk_AC = pairing.NewZr().SetFromHash([]byte(str1))

	cT1 := pairing.NewG1().MulZn(Ta1, c1).String()
	cT2 := pairing.NewG1().MulZn(Ta2, c2).String()
	Kca_str := Kca.String()

	str2 := ID_c + ID_A + cT1 + cT2 + Kca_str + Tc1_str + Tc2_str + Ta1_str + Ta2_str
	sk_CA = pairing.NewZr().SetFromHash([]byte(str2))

	fmt.Printf("sk_AC ?== sk_CA: %v \n", sk_AC.Equals(sk_CA))

	sa = pairing.NewG1().MulZn(Ta2, c1)
	sb = pairing.NewG1().MulZn(Ta1, c2)

	sa_j0 := pairing.NewG1().Add(sa, sb)
	sa_j = pairing.NewZr().SetFromHash(sa_j0.Bytes())

	tt[2] = time.Since(begin)
	fmt.Printf("the time of Medical records management is tt[2]: %v \n", tt[2])

	fmt.Printf("sa = %s\n", sa)
	fmt.Printf("sb = %s\n", sb)

}

func Med_sharing() {
	fmt.Println("----------Medical sharing 阶段----------")
	begin := time.Now()

	//doctor B 注册身份
	ID_b = "1010010422"
	Qb = pairing.NewZr().SetFromHash([]byte(ID_b))

	Sb0 := pairing.NewZr().Add(s, Qb)
	Sb1 := pairing.NewZr().Invert(Sb0)
	Sb = pairing.NewG1().MulZn(P, Sb1)

	Qb_1 = pairing.NewG1().SetFromHash([]byte(ID_b))
	Sb_1 = pairing.NewG1().MulZn(Qb_1, s)

	//patient C 生成的参数
	Kcb = pairing.NewGT().Pair(Sc_1, Qb_1)

	c1 := pairing.NewZr().Rand()
	c2 := pairing.NewZr().Rand()

	Tc1 = pairing.NewG1().MulZn(P, c1)
	Tc2 = pairing.NewG1().MulZn(P, c2)

	//doctor B 生成的参数
	b1 := pairing.NewZr().Rand()
	b2 := pairing.NewZr().Rand()

	Tb1 = pairing.NewG1().MulZn(P, b1)
	Tb2 = pairing.NewG1().MulZn(P, b2)

	Kbc = pairing.NewGT().Pair(Sb_1, Qc_1)

	fmt.Printf("Kbc ?== Kcb: %v \n", Kcb.Equals(Kbc))


	//B和C 协商秘钥
	bT1 := pairing.NewG1().MulZn(Tc1, b1).String()
	bT2 := pairing.NewG1().MulZn(Tc2, b2).String()
	Kbc_str := Kbc.String()
	Tc1_str := Tc1.String()
	Tc2_str := Tc2.String()
	Tb1_str := Tb1.String()
	Tb2_str := Tb2.String()

	str1 := ID_c + ID_A + bT1 + bT2 + Kbc_str + Tc1_str + Tc2_str + Tb1_str + Tb2_str
	sk_BC = pairing.NewZr().SetFromHash([]byte(str1))

	cT1 := pairing.NewG1().MulZn(Tb1, c1).String()
	cT2 := pairing.NewG1().MulZn(Tb2, c2).String()
	Kcb_str := Kcb.String()

	str2 := ID_c + ID_A + cT1 + cT2 + Kcb_str + Tc1_str + Tc2_str + Tb1_str + Tb2_str
	sk_CB = pairing.NewZr().SetFromHash([]byte(str2))

	fmt.Printf("sk_BC ?== sk_CB: %v \n", sk_BC.Equals(sk_CB))

	xb := pairing.NewZr().Rand()
	LB1 := pairing.NewG1().MulZn(P, xb)

	yb := pairing.NewZr().Rand()
	LB2 := pairing.NewG1().MulZn(P, yb)

	LB = pairing.NewG1().Add(LB1, LB2)

	//C 计算参数
	QA := pairing.NewZr().SetFromHash([]byte(ID_A))
	QB := pairing.NewZr().SetFromHash([]byte(ID_b))

	Ac0 := pairing.NewZr().Add(c1, QA).ThenAdd(QB)
	Ac1 := pairing.NewZr().Invert(Ac0)
	Ac = pairing.NewG1().MulZn(P, Ac1)

	CP0 := pairing.NewG1().MulZn(sa, c1)
	CP = pairing.NewG1().Mul(P, CP0)

	//B 计算参数
	b := pairing.NewZr().Rand()
	Tb := pairing.NewG1().MulZn(P, b)
	gb := pairing.NewGT().PowZn(g, b)

	RB0 := pairing.NewG1().MulZn(P, QB)
	RB1 := pairing.NewG1().Add(P_pub, RB0)
	RB := pairing.NewG1().MulZn(RB1, b)

	a_CB := pairing.NewG1().MulZn(Ac, b)

	m = "request"
	Hm := pairing.NewZr().SetFromHash([]byte(m))
	Hgb := pairing.NewZr().SetFromHash(gb.Bytes())

	N := pairing.NewZr().Add(Hgb, Hm)

	h0 := m + ID_b + P_pub.String() + LB.String() + Tb.String()
	h = pairing.NewZr().SetFromHash([]byte(h0))

	f0 := ID_b + P_pub.String() + LB.String()
	f := pairing.NewZr().SetFromHash([]byte(f0))

	dB0 := pairing.NewZr().Mul(s, f)
	dB1 := pairing.NewZr().Div(dB0, q)
	dB := pairing.NewZr().Add(yb, dB1)

	v0 := pairing.NewZr().Add(dB, xb)
	v1 := pairing.NewZr().Mul(h, v0).ThenDiv(q)
	v := pairing.NewZr().Add(b, v1)

	//hospital A 计算参数
	TQ0 := pairing.NewG1().MulZn(P, QA)
	TQ1 := pairing.NewG1().MulZn(P, QB)
	TQ := pairing.NewG1().Add(Tc1, TQ0).ThenAdd(TQ1)

	gb1 := pairing.NewGT().Pair(TQ, a_CB)

	Hgb1 := pairing.NewZr().SetFromHash(gb1.Bytes())
	m1 := pairing.NewZr().Sub(N, Hgb1)

	f10 := ID_b + P_pub.String() + LB.String()
	f1 := pairing.NewZr().SetFromHash([]byte(f10))

	Tb_10 := pairing.NewG1().MulZn(P, v)
	Tb_11 := pairing.NewG1().MulZn(P_pub, f1)
	Tb_12 := pairing.NewG1().Add(LB, Tb_11).ThenMulZn(h)
	Tb_1 := pairing.NewG1().Sub(Tb_10, Tb_12)

	h10 := m1.String() + ID_b + P_pub.String() + LB.String() + Tb_1.String()
	h1 := pairing.NewZr().SetFromHash([]byte(h10))

	fmt.Printf("h1 ?== h: %v \n", h1.Equals(h))

	//生成sk_AB参数
	K_AB := pairing.NewGT().Pair(SA_1, Qb_1)

	K_AB_C0 := pairing.NewZr().Mul(sa_j, a1)
	K_AB_C1 := pairing.NewGT().Pair(Tb, Tc1)
	K_AB_C := pairing.NewGT().PowZn(K_AB_C1, K_AB_C0)

	aTb := pairing.NewG1().MulZn(Tb, a1)
	sk_AB0 := ID_b + ID_A + aTb.String() + K_AB_C.String() + K_AB.String() + Tb.String() + Ta1.String()
	sk_AB := pairing.NewZr().SetFromHash([]byte(sk_AB0))

	//生成sk_BA参数
	K_BA := pairing.NewGT().Pair(Sb_1, QA_1)

	K_BA_C0 := pairing.NewG1().MulZn(P,sa_j).ThenMulZn(c1)
	K_BA_C1 := pairing.NewGT().Pair(Ta1, K_BA_C0)
	K_BA_C := pairing.NewGT().PowZn(K_BA_C1, b)

	bTa := pairing.NewG1().MulZn(Ta1, b)
	sk_BA0 := ID_b + ID_A + bTa.String() + K_BA_C.String() + K_BA.String() + Tb.String() + Ta1.String()
	sk_BA := pairing.NewZr().SetFromHash([]byte(sk_BA0))

	fmt.Printf("sk_AB ?== sk_BA : %v \n", sk_AB.Equals(sk_BA))

	tt[3] = time.Since(begin)
	fmt.Printf("the time of Medical sharing is tt[3]: %v \n", tt[3])

	fmt.Printf("RB = %s\n", RB)

}

// RSA签名 - 私钥
func SignatureRSA(plainText []byte, fileName string) []byte{
	//1. 打开磁盘的私钥文件
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	//2. 将私钥文件中的内容读出
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	file.Close()
	//3. 使用pem对数据解码, 得到了pem.Block结构体变量
	block, _ := pem.Decode(buf)
	//4. x509将数据解析成私钥结构体 -> 得到了私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//5. 创建一个哈希对象 -> md5/sha1 -> sha512
	// sha512.Sum512()
	myhash := sha512.New()
	//6. 给哈希对象添加数据
	myhash.Write(plainText)
	//7. 计算哈希值
	hashText := myhash.Sum(nil)
	//8. 使用rsa中的函数对散列值签名
	sigText, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashText)
	if err != nil {
		panic(err)
	}
	return sigText
}


