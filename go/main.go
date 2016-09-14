package main

import "fmt"
import "./bls"

func verifyTrue(b bool) {
	if !b {
		fmt.Println("ERR")
	}
}
func testRecoverSecretKey() {
	fmt.Println("testRecoverSecretKey")
	k := 1000
	sec := bls.NewSecretKey()
	sec.Init()
	fmt.Println("sec =", sec)

	// make master secret key
	msk := sec.GetMasterSecretKey(k)

	n := k
	secVec := make([]*bls.SecretKey, n)
	idVec := make([]*bls.Id, n)
	for i := 0; i < n; i++ {
		idVec[i] = bls.NewId()
		idVec[i].Set([]uint64{1, 2, 3, uint64(i)})
		secVec[i] = bls.NewSecretKey()
		secVec[i].Set(msk, idVec[i])
	}
	// recover sec2 from secVec and idVec
	sec2 := bls.NewSecretKey()
	sec2.Recover(secVec, idVec)
	fmt.Println("sec2=", sec2)
}

func testSign() {
	m := "testSign"
	fmt.Println(m)

	sec0 := bls.NewSecretKey()
	sec0.Init()
	pub0 := sec0.GetPublicKey()
	s0 := sec0.Sign(m)
	verifyTrue(s0.Verify(pub0, m))

	k := 3
	msk := sec0.GetMasterSecretKey(k)
	mpk := bls.GetMasterPublicKey(msk)

	idTbl := []uint64{3, 5, 193, 22, 15}
	n := len(idTbl)

	secVec := make([]*bls.SecretKey, n)
	pubVec := make([]*bls.PublicKey, n)
	signVec := make([]*bls.Sign, n)
	idVec := make([]*bls.Id, n)

	for i := 0; i < n; i++ {
		idVec[i] = bls.NewId()
		idVec[i].Set([]uint64{idTbl[i], 0, 0, 0})
		fmt.Printf("idVec[%d]=%s\n", i, idVec[i].String())

		secVec[i] = bls.NewSecretKey()
		secVec[i].Set(msk, idVec[i])

		pubVec[i] = bls.NewPublicKey()
		pubVec[i].Set(mpk, idVec[i])
		fmt.Printf("pubVec[%d]=%s\n", i, pubVec[i].String())

		verifyTrue(pubVec[i].String() == secVec[i].GetPublicKey().String())

		signVec[i] = secVec[i].Sign(m)
		verifyTrue(signVec[i].Verify(pubVec[i], m))
	}
	sec1 := bls.NewSecretKey()
	sec1.Recover(secVec, idVec)
	verifyTrue(sec0.String() == sec1.String())
	pub1 := bls.NewPublicKey()
	pub1.Recover(pubVec, idVec)
	verifyTrue(pub0.String() == pub1.String())
	s1 := bls.NewSign()
	s1.Recover(signVec, idVec)
	verifyTrue(s0.String() == s1.String())
}

func testAdd() {
	fmt.Println("testAdd")
	sec1 := bls.NewSecretKey()
	sec2 := bls.NewSecretKey()
	sec1.Init()
	sec2.Init()

	pub1 := sec1.GetPublicKey()
	pub2 := sec2.GetPublicKey()

	m := "test test"
	sign1 := sec1.Sign(m)
	sign2 := sec2.Sign(m)

	fmt.Println("sign1    :", sign1)
	sign1.Add(sign2)
	fmt.Println("sign1 add:", sign1)
	pub1.Add(pub2)
	verifyTrue(sign1.Verify(pub1, m))
}

func testPop() {
	fmt.Println("testPop")
	sec := bls.NewSecretKey()
	sec.Init()
	pop := sec.GetPop()
	verifyTrue(pop.VerifyPop(sec.GetPublicKey()))
	sec.Init()
	verifyTrue(!pop.VerifyPop(sec.GetPublicKey()))
}
func main() {
	fmt.Println("init")
	bls.Init()
	{
		id := bls.NewId()
		id.Set([]uint64{4, 3, 2, 1})
		fmt.Println("id :", id)
		id2 := bls.NewId()
		id2.SetStr(id.String())
		fmt.Println("id2:", id2)
	}
	{
		sec := bls.NewSecretKey()
		sec.SetArray([]uint64{1, 2, 3, 4})
		fmt.Println("sec=", sec)
	}

	fmt.Println("create secret key")
	m := "this is a bls sample for go"
	sec := bls.NewSecretKey()
	sec.Init()
	fmt.Println("sec:", sec)
	fmt.Println("create public key")
	pub := sec.GetPublicKey()
	fmt.Println("pub:", pub)
	sign := sec.Sign(m)
	fmt.Println("sign:", sign)
	verifyTrue(sign.Verify(pub, m))

	// How to make array of SecretKey
	{
		sec := make([]bls.SecretKey, 3)
		for i := 0; i < len(sec); i++ {
			sec[i] = *bls.NewSecretKey()
			sec[i].Init()
			fmt.Println("sec=", sec[i].String())
		}
	}
	testRecoverSecretKey()
	testAdd()
	testSign()
	testPop()
}
